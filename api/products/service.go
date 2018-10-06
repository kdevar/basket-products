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
	queries := []elastic.Query{}
	productIdTermQuery := elastic.NewTermQuery(_const.PRODUCTIDFIELD, filter.ProductId)
	queries = append(queries, productIdTermQuery)

	priceStatusFilter := elastic.NewTermsQuery(_const.PRICESTATUSFIELD, "0")
	productStatusFilter := elastic.NewTermsQuery(_const.PRODUCTSTATUSFIELD, "0")

	queries = append(queries, priceStatusFilter, productStatusFilter)

	if len(filter.ChainId) > 0 {
		s := util.ConvertStringToInterface(filter.ChainId)
		chainIdTermQuery := elastic.NewTermsQuery("chainId", s...)
		queries = append(queries, chainIdTermQuery)
	}

	gq := elastic.NewBoolQuery().Filter(queries...)

	baseAggregation := elastic.NewTermsAggregation().Field(_const.CHAINIDFIELD)

	ma := NewEstimateTermsAgg(_const.METROAREAIDFIELD, filter.MetroAreaId).Generate()

	c := NewEstimateTermsAgg(_const.CITYIDFIELD, filter.CityId).Generate()

	z := NewEstimateTermsAgg(_const.ZIPIDFIELD, filter.ZipCodeId).Generate()

	a := elastic.
		NewGeoDistanceAggregation().
		Field(_const.LOCATIONFIELD).
		Point(filter.GetLatLongString()).
		Unit("miles").
		AddRangeWithKey(string(_const.CITY), 0, 50).
		AddRangeWithKey(string(_const.METRO), 51, 100).
		AddRangeWithKey(string(_const.NATIONALMILES), 101, 4000)

	s := elastic.NewMinAggregation().Field(_const.FINALPRICEFIELD)
	m := elastic.NewMaxAggregation().Field(_const.FINALPRICEFIELD)

	a.SubAggregation(_const.MINLABEL, s)
	a.SubAggregation(_const.MAXLABEL, m)

	baseAggregation.SubAggregation(_const.GEOGRAPHICLABEL, a)
	baseAggregation.SubAggregation(_const.METROLABEL, ma)
	baseAggregation.SubAggregation(_const.CITYLABEL, c)
	baseAggregation.SubAggregation(_const.ZIPLABEL, z)

	searchQuery := svc.elasticClient.
		Search(_const.PROUCTPRICEINDEX).
		Query(gq).
		Aggregation(_const.CHAINLABEL, baseAggregation)

	searchResult, err := searchQuery.Do(ctx)

	if err != nil {
		errors.ServerError(err)
	}

	groups, _ := searchResult.Aggregations.Terms(_const.CHAINLABEL)
	responses := make(EstimatedPriceResponse)
	for _, b := range groups.Buckets {
		chainId := ChainId(b.KeyNumber)
		estimates := make(EstimatedPriceItem)

		cityEstimate := &ProductEstimate{
			Etype: _const.CITY,
		}
		cityBucket, _ := b.Aggregations.Terms(_const.CITYLABEL)
		cityEstimate.transformFromTermsBucket(cityBucket)

		estimates[_const.CITY] = cityEstimate

		metroEstimate := &ProductEstimate{
			Etype: _const.METRO,
		}
		metroBucket, _ := b.Aggregations.Terms(_const.METROLABEL)
		metroEstimate.transformFromTermsBucket(metroBucket)

		estimates[_const.METRO] = metroEstimate

		zipEstimate := &ProductEstimate{
			Etype: "ZIP",
		}
		zipBucket, _ := b.Aggregations.Terms(_const.ZIPLABEL)
		zipEstimate.transformFromTermsBucket(zipBucket)

		estimates[_const.ZIP] = zipEstimate

		x, _ := b.Aggregations.GeoDistance(_const.GEOGRAPHICLABEL)

		for _, k := range x.Buckets {
			max, _ := k.Aggregations.Max(_const.MAXLABEL)
			min, _ := k.Aggregations.Min(_const.MINLABEL)
			if k.Key == string(_const.CITY) {
				estimates[_const.FITYMILES] = &ProductEstimate{
					Etype: _const.FITYMILES,
					Min:   min.Value,
					Max:   max.Value}
			}
			if k.Key == string(_const.METRO) {
				estimates[_const.HUNDREDMILES] = &ProductEstimate{
					Etype: _const.HUNDREDMILES,
					Min:   min.Value, Max: max.Value}
			}
			if k.Key == string(_const.NATIONALMILES) {
				estimates[_const.NATIONALMILES] = &ProductEstimate{
					Etype: _const.NATIONALMILES,
					Min:   min.Value,
					Max:   max.Value}
			}
		}

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

	fq.Query(bq).AddScoreFunc(businessValueScore).AddScoreFunc(popularityScore).AddScoreFunc(totalPricesScore)

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
