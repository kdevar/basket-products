package products

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kdevar/basket-products/const"
	"github.com/kdevar/basket-products/errors"
	"github.com/kdevar/basket-products/util"
	"github.com/olivere/elastic"
)

type productServiceImpl struct {
	elasticClient *elastic.Client
}

func (svc *productServiceImpl) GetEstimatedProductPrices(filter EstimatedPriceFilter) (EstimatedPriceResponse, *errors.ApiError) {
	ctx := context.Background()
	query := NewEstimatedQuery().
		AddProductIdFilter(filter.ProductId).
		AddPriceStatusFilter().
		AddProductStatusFilter()

	if filter.ChainId != nil {
		query.AddChainIdFilters(filter.ChainId)
	}

	agg := NewEstimationAggregation().
		AddMetroAreaSubAggregation(filter.MetroAreaId).
		AddCityIdSubAggregation(filter.CityId).
		AddZipIdSubAggregation(filter.ZipCodeId).
		AddGeographicRangeSubAggregation(filter.Location)

	searchQuery := svc.elasticClient.
		Search(_const.PROUCTPRICEINDEX).
		Query(query.GetQuery()).
		Aggregation(_const.CHAINLABEL, agg.GetAggregation())

	searchResult, err := searchQuery.Do(ctx)

	if err != nil {
		return nil, errors.ServerError(err)
	}

	groups, _ := searchResult.Aggregations.Terms(_const.CHAINLABEL)

	chainGroups := NewEstimationAggregationResultIterator(groups)
	responses := make(EstimatedPriceResponse)
	for chainGroups.Next(){
		estimates := make(EstimatedPriceItem)
		chainGroup := chainGroups.Value()
		chainId := ChainId(chainGroup.Result.KeyNumber)
		estimates[_const.CITY] = chainGroup.GetCityEstimate()
		estimates[_const.METRO] = chainGroup.GetMetroEstimate()
		estimates[_const.ZIP] = chainGroup.GetZipEstimate()
		estimates[_const.FIFTYMILES] = chainGroup.GetRangeEstimate(_const.FIFTYMILES)
		estimates[_const.HUNDREDMILES] = chainGroup.GetRangeEstimate(_const.HUNDREDMILES)
		estimates[_const.NATIONALMILES] = chainGroup.GetRangeEstimate(_const.NATIONALMILES)
		responses[chainId] = estimates
	}
	return responses, nil
}

func (svc *productServiceImpl) SearchProducts(filter SearchFilter) ([]Product, *errors.ApiError) {
	ctx := context.Background()
	filters := []elastic.Query{}
	fq := elastic.NewFunctionScoreQuery()

	priceStatusFilter := elastic.NewTermsQuery(_const.PRICESTATUSFIELD, "0")
	productStatusFilter := elastic.NewTermsQuery(_const.PRODUCTSTATUSFIELD, "0")

	keywordMultiMatch := elastic.NewMultiMatchQuery(
		filter.keyword,
		"productDesc.raw^5",
		"brandDesc.raw^3",
		"typeDesc.raw^10",
		"fullProductName^10",
		"tags")

	filters = append(filters, productStatusFilter)
	filters = append(filters, priceStatusFilter)

	if filter.location != nil {
		geoLocationFilter := elastic.NewGeoDistanceQuery(_const.LOCATIONFIELD).GeoPoint(filter.location).Distance("25mi")
		filters = append(filters, geoLocationFilter)
	}

	if filter.productIds != nil {
		s := util.ConvertIntToInterface(filter.productIds)
		productIdFilter := elastic.NewTermsQuery(_const.PRODUCTIDFIELD, s...)
		filters = append(filters, productIdFilter)
	}

	if filter.categoryId != nil {
		categoryIdFilter := elastic.NewTermsQuery(_const.CATEGORYIDFIELD, *filter.categoryId)
		filters = append(filters, categoryIdFilter)
	}

	if filter.brandId != nil {
		brandIdFilter := elastic.NewTermsQuery(_const.BRANDIDFIELD, *filter.brandId)
		filters = append(filters, brandIdFilter)
	}

	if filter.typeId != nil {
		typeIdFilter := elastic.NewTermsQuery(_const.TYPEIDFIELD, *filter.typeId)
		filters = append(filters, typeIdFilter)
	}

	businessValueScore := elastic.
		NewFieldValueFactorFunction().
		Field(_const.BUSINESSVALUEFIELD).
		Missing(1.0)
	popularityScore := elastic.
		NewFieldValueFactorFunction().
		Field(_const.POPULARITYFIELD).
		Factor(0.001).
		Missing(1.0)
	totalPricesScore := elastic.
		NewFieldValueFactorFunction().
		Field(_const.TOTALPRICESCOREFIELD).
		Factor(0.0001).
		Missing(1.0)

	bq := elastic.NewBoolQuery().Should(keywordMultiMatch).Filter(filters...)

	docValueFields := []string{
		_const.PRODUCTIDFIELD,
		_const.STOREIDFIELD,
		_const.LISTPRICEPERONEFIELD,
		_const.LISTPRICEFIELD,
		_const.LISTQUANTITYFIELD,
		_const.LISTPRICEUSERIDFIELD,
		_const.SALEPRICEPERONEFIELD,
		_const.SALEPRICETYPEID,
		_const.SALEVALUE1FIELD,
		_const.SALEVALUE2FIELD,
		_const.SALEVALUE3FIELD,
		_const.SALEPRICEUSERIDFIELD,
		_const.SALEENDDATEFIELD,
		_const.COUPONPRICEPERONEFIELD,
		_const.COUPONPRICETYPEIDFIELD,
		_const.COUPONVALUE1FIELD,
		_const.COUPONVALUE2FIELD,
		_const.COUPONVALUE3FIELD,
		_const.COUPONURIFIELD,
		_const.COUPONENDDATEFIELD,
		_const.UNITPRICEPERONEFIELD,
		_const.UNIDESCFIELD,
		_const.FINALPRICEFIELD,
	}

	finalPriceInnerHit := elastic.NewInnerHit().
		Name("MINIMAL").
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include(_const.PRODUCTIDFIELD)).
		DocvalueFields(docValueFields...).
		Sort(_const.FINALPRICEFIELD, true).
		Size(1)

	listPricePerOneInnerHit := elastic.NewInnerHit().
		Name("MAXIMUM").
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include(_const.PRODUCTIDFIELD)).
		DocvalueFields(docValueFields...).
		Sort(_const.LISTPRICEPERONEFIELD, false).
		Size(1)

	collapseBuilder := elastic.
		NewCollapseBuilder(_const.PRODUCTIDFIELD).
		InnerHit(finalPriceInnerHit).
		InnerHit(listPricePerOneInnerHit)

	searchQuery := svc.elasticClient.
		Search(_const.PROUCTPRICEINDEX).
		Query(fq.
			Query(bq).
			AddScoreFunc(businessValueScore).
			AddScoreFunc(popularityScore).
			AddScoreFunc(totalPricesScore).
			ScoreMode("multiply").
			BoostMode("multiply")).
		Collapse(collapseBuilder)

	if filter.from != nil && filter.to != nil {
		searchQuery.From(*filter.from).Size(*filter.to)
	}

	searchResult, searchErr := searchQuery.Do(ctx)

	if searchErr != nil {
		return nil, errors.ServerError(searchErr)
	}

	var products []Product

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var t Product
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				fmt.Println(err)
			}
			products = append(products, t)
		}
	}
	return products, nil
}

func (svc *productServiceImpl) GetLiveProductPrices(filter LivePriceFilter) ([]Product, *errors.ApiError) {
	ctx := context.Background()

	fq := elastic.NewFunctionScoreQuery()

	productIdQuery := elastic.NewTermQuery(_const.PRODUCTIDFIELD, filter.productId)
	priceStatusFilter := elastic.NewTermsQuery(_const.PRICESTATUSFIELD, "0")
	productStatusFilter := elastic.NewTermsQuery(_const.PRODUCTSTATUSFIELD, "0")

	bq := elastic.NewBoolQuery().Filter(productIdQuery).Filter(priceStatusFilter, productStatusFilter)

	if len(filter.storeIds) > 0 {
		s := util.ConvertStringToInterface(filter.storeIds)
		storeIdFilter := elastic.NewTermsQuery("storeId", s...)
		bq.Filter(storeIdFilter)
	}

	businessValueScore := elastic.
		NewFieldValueFactorFunction().
		Field(_const.BUSINESSVALUEFIELD).
		Missing(1.0)
	popularityScore := elastic.NewFieldValueFactorFunction().
		Field(_const.POPULARITYFIELD).
		Factor(0.001).
		Missing(1.0)
	totalPricesScore := elastic.
		NewFieldValueFactorFunction().
		Field(_const.TOTALPRICESCOREFIELD).
		Factor(0.0001).
		Missing(1.0)

	if filter.location != nil {
		geoLocationFilter := elastic.NewGeoDistanceQuery(_const.LOCATIONFIELD).GeoPoint(filter.location).Distance("25mi")
		bq.Filter(geoLocationFilter)
	}

	fq.Query(bq).
		AddScoreFunc(businessValueScore).
		AddScoreFunc(popularityScore).
		AddScoreFunc(totalPricesScore)

	searchResult, err := svc.elasticClient.Search(_const.PROUCTPRICEINDEX).Query(fq).Do(ctx)

	if err != nil {
		errors.ServerError(err)
	}
	var products []Product

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var t Product
			for _, innerHit := range hit.InnerHits {
				var ih map[string]interface{}
				for _, h := range innerHit.Hits.Hits {
					fmt.Println(json.Unmarshal(*h.Source, &ih))
				}
			}
			err := json.Unmarshal(*hit.Source, &t)

			if err != nil {
				fmt.Println(err)
			}
			products = append(products, t)
		}
	} else {
		fmt.Print("Found no products\n")
	}

	return products, nil

}

type ProductEstimate struct {
	Etype _const.EstimateType `json:"Etype"`
	Min   *float64            `json:"Min"`
	Max   *float64            `json:"Max"`
}

func (e *ProductEstimate) transformFromTermsBucket(bucket *elastic.AggregationBucketKeyItems) {
	if len(bucket.Buckets) == 1 {
		max, ok := bucket.Buckets[0].Aggregations.Max(_const.MAXFINALPRICELABEL)
		if ok {
			e.Max = max.Value
		}
		min, ok := bucket.Buckets[0].Aggregations.Min(_const.MINFINALPRICELABEL)
		if ok {
			e.Min = min.Value
		}
	}
}

type EstimateTermsAgg struct {
	fieldName    string
	inclusionVal string
}

func (a *EstimateTermsAgg) Generate() *elastic.TermsAggregation {
	return elastic.
		NewTermsAggregation().
		Field(a.fieldName).
		Include(a.inclusionVal).
		SubAggregation(_const.MINFINALPRICELABEL, elastic.NewMinAggregation().Field(_const.FINALPRICEFIELD)).
		SubAggregation(_const.MAXFINALPRICELABEL, elastic.NewMaxAggregation().Field(_const.FINALPRICEFIELD)).
		SubAggregation(_const.MINLISTLABEL, elastic.NewMinAggregation().Field(_const.LISTPRICEFIELD)).
		SubAggregation(_const.MAXLISTLABEL, elastic.NewMaxAggregation().Field(_const.LISTPRICEFIELD))
}

func NewEstimateTermsAgg(fieldName string, inclusionVal string) *EstimateTermsAgg {
	return &EstimateTermsAgg{
		fieldName:    fieldName,
		inclusionVal: inclusionVal,
	}
}

type ChainId json.Number
type EstimatedPriceItem map[_const.EstimateType]*ProductEstimate
type EstimatedPriceResponse map[ChainId]EstimatedPriceItem
