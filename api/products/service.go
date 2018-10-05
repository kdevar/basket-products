package products

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kdevar/basket-products/api/errors"
	"github.com/kdevar/basket-products/util"
	"github.com/olivere/elastic"
)

type productServiceImpl struct {
	elasticClient *elastic.Client
}

func (svc *productServiceImpl) GetEstimatedProductPrices(filter EstimatedPriceFilter) (EstimatedPriceResponse, *errors.ApiError) {
	ctx := context.Background()
	queries := []elastic.Query{}
	productIdTermQuery := elastic.NewTermQuery(PRODUCTIDFIELD, filter.ProductId)
	queries = append(queries, productIdTermQuery)

	priceStatusFilter := elastic.NewTermsQuery(PRICESTATUSFIELD, "0")
	productStatusFilter := elastic.NewTermsQuery(PRODUCTSTATUSFIELD, "0")

	queries = append(queries, priceStatusFilter, productStatusFilter)

	if len(filter.ChainId) > 0 {
		s := util.ConvertSToInterface(filter.ChainId)
		chainIdTermQuery := elastic.NewTermsQuery("chainId", s...)
		queries = append(queries, chainIdTermQuery)
	}

	gq := elastic.NewBoolQuery().Filter(queries...)

	baseAggregation := elastic.NewTermsAggregation().Field(CHAINIDFIELD)

	ma := NewEstimateTermsAgg(METROAREAIDFIELD, filter.MetroAreaId).Generate()

	c := NewEstimateTermsAgg(CITYIDFIELD, filter.CityId).Generate()

	z := NewEstimateTermsAgg(ZIPIDFIELD, filter.ZipCodeId).Generate()

	a := elastic.
		NewGeoDistanceAggregation().
		Field(LOCATIONFIELD).
		Point(filter.GetLatLongString()).
		Unit("miles").
		AddRangeWithKey(string(CITY), 0, 50).
		AddRangeWithKey(string(METRO), 51, 100).
		AddRangeWithKey(string(NATIONALMILES), 101, 4000)

	s := elastic.NewMinAggregation().Field(FINALPRICEFIELD)
	m := elastic.NewMaxAggregation().Field(FINALPRICEFIELD)

	a.SubAggregation(MINLABEL, s)
	a.SubAggregation(MAXLABEL, m)

	baseAggregation.SubAggregation(GEOGRAPHICLABEL, a)
	baseAggregation.SubAggregation(METROLABEL, ma)
	baseAggregation.SubAggregation(CITYLABEL, c)
	baseAggregation.SubAggregation(ZIPLABEL, z)

	searchQuery := svc.elasticClient.
		Search(PROUCTPRICEINDEX).
		Query(gq).
		Aggregation(CHAINLABEL, baseAggregation)

	searchResult, err := searchQuery.Do(ctx)

	if err != nil {
		errors.ServerError(err)
	}

	groups, _ := searchResult.Aggregations.Terms(CHAINLABEL)
	responses := make(EstimatedPriceResponse)
	for _, b := range groups.Buckets {
		chainId := ChainId(b.KeyNumber)
		estimates := make(EstimatedPriceItem)

		cityEstimate := &ProductEstimate{
			Etype: CITY,
		}
		cityBucket, _ := b.Aggregations.Terms(CITYLABEL)
		cityEstimate.transformFromTermsBucket(cityBucket)

		estimates[CITY] = cityEstimate

		metroEstimate := &ProductEstimate{
			Etype: METRO,
		}
		metroBucket, _ := b.Aggregations.Terms(METROLABEL)
		metroEstimate.transformFromTermsBucket(metroBucket)

		estimates[METRO] = metroEstimate

		zipEstimate := &ProductEstimate{
			Etype: "ZIP",
		}
		zipBucket, _ := b.Aggregations.Terms(ZIPLABEL)
		zipEstimate.transformFromTermsBucket(zipBucket)

		estimates[ZIP] = zipEstimate

		x, _ := b.Aggregations.GeoDistance(GEOGRAPHICLABEL)

		for _, k := range x.Buckets {
			max, _ := k.Aggregations.Max(MAXLABEL)
			min, _ := k.Aggregations.Min(MINLABEL)
			if k.Key == string(CITY) {
				estimates[FITYMILES] = &ProductEstimate{
					Etype: FITYMILES,
					Min:   min.Value,
					Max:   max.Value}
			}
			if k.Key == string(METRO) {
				estimates[HUNDREDMILES] = &ProductEstimate{
					Etype: HUNDREDMILES,
					Min:   min.Value, Max: max.Value}
			}
			if k.Key == string(NATIONALMILES) {
				estimates[NATIONALMILES] = &ProductEstimate{
					Etype: NATIONALMILES,
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

	priceStatusFilter := elastic.NewTermsQuery(PRICESTATUSFIELD, "0")
	productStatusFilter := elastic.NewTermsQuery(PRODUCTSTATUSFIELD, "0")

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
		geoLocationFilter := elastic.NewGeoDistanceQuery(LOCATIONFIELD).GeoPoint(filter.location).Distance("25mi")
		filters = append(filters, geoLocationFilter)
	}

	if filter.productIds != nil {
		s := util.ConvertIToInterface(filter.productIds)
		productIdFilter := elastic.NewTermsQuery(PRODUCTIDFIELD, s...)
		filters = append(filters, productIdFilter)
	}

	if filter.categoryId != nil {
		categoryIdFilter := elastic.NewTermsQuery(CATEGORYIDFIELD, *filter.categoryId)
		filters = append(filters, categoryIdFilter)
	}

	if filter.brandId != nil {
		brandIdFilter := elastic.NewTermsQuery(BRANDIDFIELD, *filter.brandId)
		filters = append(filters, brandIdFilter)
	}

	if filter.typeId != nil {
		typeIdFilter := elastic.NewTermsQuery(TYPEIDFIELD, *filter.typeId)
		filters = append(filters, typeIdFilter)
	}

	businessValueScore := elastic.
		NewFieldValueFactorFunction().
		Field(BUSINESSVALUEFIELD).
		Missing(1.0)
	popularityScore := elastic.
		NewFieldValueFactorFunction().
		Field(POPULARITYFIELD).
		Factor(0.001).
		Missing(1.0)
	totalPricesScore := elastic.
		NewFieldValueFactorFunction().
		Field(TOTALPRICESCOREFIELD).
		Factor(0.0001).
		Missing(1.0)

	bq := elastic.NewBoolQuery().Should(keywordMultiMatch).Filter(filters...)

	docValueFields := []string{
		PRODUCTIDFIELD,
		STOREIDFIELD,
		LISTPRICEPERONEFIELD,
		LISTPRICEFIELD,
		LISTQUANTITYFIELD,
		LISTPRICEUSERIDFIELD,
		SALEPRICEPERONEFIELD,
		SALEPRICETYPEID,
		SALEVALUE1FIELD,
		SALEVALUE2FIELD,
		SALEVALUE3FIELD,
		SALEPRICEUSERIDFIELD,
		SALEENDDATEFIELD,
		COUPONPRICEPERONEFIELD,
		COUPONPRICETYPEIDFIELD,
		COUPONVALUE1FIELD,
		COUPONVALUE2FIELD,
		COUPONVALUE3FIELD,
		COUPONURIFIELD,
		COUPONENDDATEFIELD,
		UNITPRICEPERONEFIELD,
		UNIDESCFIELD,
		FINALPRICEFIELD,
	}

	finalPriceInnerHit := elastic.NewInnerHit().
		Name("MINIMAL").
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include(PRODUCTIDFIELD)).
		DocvalueFields(docValueFields...).
		Sort(FINALPRICEFIELD, true).
		Size(1)

	listPricePerOneInnerHit := elastic.NewInnerHit().
		Name("MAXIMUM").
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include(PRODUCTIDFIELD)).
		DocvalueFields(docValueFields...).
		Sort(LISTPRICEPERONEFIELD, false).
		Size(1)

	collapseBuilder := elastic.
		NewCollapseBuilder(PRODUCTIDFIELD).
		InnerHit(finalPriceInnerHit).
		InnerHit(listPricePerOneInnerHit)

	searchQuery := svc.elasticClient.
		Search(PROUCTPRICEINDEX).
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

	productIdQuery := elastic.NewTermQuery(PRODUCTIDFIELD, filter.productId)
	priceStatusFilter := elastic.NewTermsQuery(PRICESTATUSFIELD, "0")
	productStatusFilter := elastic.NewTermsQuery(PRODUCTSTATUSFIELD, "0")

	bq := elastic.NewBoolQuery().Filter(productIdQuery).Filter(priceStatusFilter, productStatusFilter)

	if len(filter.storeIds) > 0 {
		s := util.ConvertSToInterface(filter.storeIds)
		storeIdFilter := elastic.NewTermsQuery("storeId", s...)
		bq.Filter(storeIdFilter)
	}

	businessValueScore := elastic.
		NewFieldValueFactorFunction().
		Field(BUSINESSVALUEFIELD).
		Missing(1.0)
	popularityScore := elastic.NewFieldValueFactorFunction().
		Field(POPULARITYFIELD).
		Factor(0.001).
		Missing(1.0)
	totalPricesScore := elastic.
		NewFieldValueFactorFunction().
		Field(TOTALPRICESCOREFIELD).
		Factor(0.0001).
		Missing(1.0)

	if filter.location != nil {
		geoLocationFilter := elastic.NewGeoDistanceQuery(LOCATIONFIELD).GeoPoint(filter.location).Distance("25mi")
		bq.Filter(geoLocationFilter)
	}

	fq.Query(bq).AddScoreFunc(businessValueScore).AddScoreFunc(popularityScore).AddScoreFunc(totalPricesScore)

	searchResult, err := svc.elasticClient.Search(PROUCTPRICEINDEX).Query(fq).Do(ctx)

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
	Etype EstimateType `json:"Etype"`
	Min   *float64     `json:"Min"`
	Max   *float64     `json:"Max"`
}

func (e *ProductEstimate) transformFromTermsBucket(bucket *elastic.AggregationBucketKeyItems) {
	if len(bucket.Buckets) == 1 {
		max, ok := bucket.Buckets[0].Aggregations.Max(MAXFINALPRICELABEL)
		if ok {
			e.Max = max.Value
		}
		min, ok := bucket.Buckets[0].Aggregations.Min(MINFINALPRICELABEL)
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
		SubAggregation(MINFINALPRICELABEL, elastic.NewMinAggregation().Field(FINALPRICEFIELD)).
		SubAggregation(MAXFINALPRICELABEL, elastic.NewMaxAggregation().Field(FINALPRICEFIELD)).
		SubAggregation(MINLISTLABEL, elastic.NewMinAggregation().Field(LISTPRICEFIELD)).
		SubAggregation(MAXLISTLABEL, elastic.NewMaxAggregation().Field(LISTPRICEFIELD))
}

func NewEstimateTermsAgg(fieldName string, inclusionVal string) *EstimateTermsAgg {
	return &EstimateTermsAgg{
		fieldName:    fieldName,
		inclusionVal: inclusionVal,
	}
}

type ChainId json.Number
type EstimatedPriceItem map[EstimateType]*ProductEstimate
type EstimatedPriceResponse map[ChainId]EstimatedPriceItem
