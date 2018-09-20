package products

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kdevar/basket-products/api/errors"
	"github.com/kdevar/basket-products/config"
	"github.com/olivere/elastic"
	"strconv"
)

var Service *productService

func init() {
	Service = &productService{}
}

type ProductFilter struct {
	keyword      string
	location     *elastic.GeoPoint
	productIds   []int
	categoryId   *string
	brandId      *string
	typeId       *string
	categoryDesc *string
}

type ProductPriceFilter struct {
	productId string
	location  *elastic.GeoPoint
	storeIds  []string
}

type productService struct{}

type ProductEstimateRequest struct {
	ProductId       string
	ChainId         []string
	MetroAreaId     string
	CityId          string
	ZipCodeId       string
	Location        *elastic.GeoPoint
	IncludeEstimate bool
	IncludeDetails  bool
}

type EstimateType string
const (
	ZIP EstimateType = "ZIP"
	CITY EstimateType = "CITY"
	METRO EstimateType = "METRO"
	FITYMILES EstimateType = "FIFTYMILE"
	HUNDREDMILES EstimateType = "HUNDREDMILES"
	NATIONALMILES EstimateType = "NATIONALMILES"
)

type ProductEstimate struct {
	Etype EstimateType  `json:"Etype"`
	Min   *float64 `json:"Min"`
	Max   *float64 `json:"Max"`
}

func (r *ProductEstimateRequest) GetLatLonString() string {
	return strconv.FormatFloat(r.Location.Lat, 'f', 6, 64) + "," + strconv.FormatFloat(r.Location.Lon, 'f', 6, 64)
}

func (svc *productService) GetProductEstimate(request ProductEstimateRequest) map[json.Number]map[EstimateType]*ProductEstimate {
	ctx := context.Background()

	queries := []elastic.Query{}

	productIdTermQuery := elastic.NewTermQuery("productId", request.ProductId)
	queries = append(queries, productIdTermQuery)


	priceStatusFilter := elastic.NewTermsQuery("priceStatus", "0")
	productStatusFilter := elastic.NewTermsQuery("productStatus", "0")

	queries = append(queries, priceStatusFilter, productStatusFilter)

	if len(request.ChainId) > 0 {
		s := make([]interface{}, len(request.ChainId))
		for i, v := range request.ChainId {
			s[i] = v
		}

		chainIdTermQuery := elastic.NewTermsQuery("chainId", s...)
		queries = append(queries, chainIdTermQuery)
	}

	gq := elastic.NewBoolQuery().Filter(queries...)

	baseAggregation := elastic.NewTermsAggregation().Field("chainId")

	ma := elastic.
		NewTermsAggregation().
		Field("metroAreaId").
		Include(request.MetroAreaId).
		SubAggregation("minfinal", elastic.NewMinAggregation().Field("finalPrice")).
		SubAggregation("maxfinal", elastic.NewMaxAggregation().Field("finalPrice")).
		SubAggregation("minlist", elastic.NewMinAggregation().Field("listPrice")).
		SubAggregation("maxlist", elastic.NewMaxAggregation().Field("listPrice"))

	c := elastic.
		NewTermsAggregation().
		Field("cityId").
		Include(request.CityId).
		SubAggregation("minfinal", elastic.NewMinAggregation().Field("finalPrice")).
		SubAggregation("maxfinal", elastic.NewMaxAggregation().Field("finalPrice")).
		SubAggregation("minlist", elastic.NewMinAggregation().Field("listPrice")).
		SubAggregation("maxlist", elastic.NewMaxAggregation().Field("listPrice"))

	z := elastic.
		NewTermsAggregation().
		Field("zipCodeId").
		Include(request.ZipCodeId).
		SubAggregation("minfinal", elastic.NewMinAggregation().Field("finalPrice")).
		SubAggregation("maxfinal", elastic.NewMaxAggregation().Field("finalPrice")).
		SubAggregation("minlist", elastic.NewMinAggregation().Field("listPrice")).
		SubAggregation("maxlist", elastic.NewMaxAggregation().Field("listPrice"))

	a := elastic.
		NewGeoDistanceAggregation().
		Field("location").
		Point(request.GetLatLonString()).
		Unit("miles").
		AddRangeWithKey("city", 0, 50).
		AddRangeWithKey("metro", 51, 100).
		AddRangeWithKey("national", 101, 4000)

	s := elastic.NewMinAggregation().Field("finalPrice")
	m := elastic.NewMaxAggregation().Field("finalPrice")

	a.SubAggregation("minprice", s)
	a.SubAggregation("maxprice", m)

	baseAggregation.SubAggregation("geographic-range-estimates", a)
	baseAggregation.SubAggregation("metro-estimate", ma)
	baseAggregation.SubAggregation("city-estimate", c)
	baseAggregation.SubAggregation("zipcode-estimate", z)

	searchQuery := config.
		ElasticClient.
		Search("prices").
		Query(gq).
		Aggregation("chain-groups", baseAggregation)

	searchResult, err := searchQuery.Do(ctx)

	if err != nil {
		fmt.Println(err)
	}

	groups, _ := searchResult.Aggregations.Terms("chain-groups")
	responses := make(map[json.Number]map[EstimateType]*ProductEstimate)
	for _, b := range groups.Buckets {

		estimates := make(map[EstimateType]*ProductEstimate)

		cityEstimate := &ProductEstimate{
			Etype: CITY,
		}
		cityEstimateAgg, _ := b.Aggregations.Terms("city-estimate")

		if len(cityEstimateAgg.Buckets) == 1 {
			max, ok := cityEstimateAgg.Buckets[0].Aggregations.Max("maxfinal")
			if ok {
				cityEstimate.Max = max.Value
			}
			min, ok := cityEstimateAgg.Buckets[0].Aggregations.Min("minfinal")

			if ok {
				cityEstimate.Min = min.Value
			}
		}

		estimates[CITY] =  cityEstimate

		metroEstimate := &ProductEstimate{
			Etype: METRO,
		}
		metroEstimateAgg, _ := b.Aggregations.Terms("metro-estimate")

		if len(metroEstimateAgg.Buckets) == 1 {
			max, ok := metroEstimateAgg.Buckets[0].Aggregations.Max("maxfinal")
			if ok {
				metroEstimate.Max = max.Value
			}
			min, ok := metroEstimateAgg.Buckets[0].Aggregations.Min("minfinal")

			if ok {
				metroEstimate.Min = min.Value
			}
		}

		estimates[METRO] =  metroEstimate

		zipEstimate := &ProductEstimate{
			Etype: "ZIP",
		}
		zipEstimateAgg, _ := b.Aggregations.Terms("zipcode-estimate")

		if len(zipEstimateAgg.Buckets) == 1 {
			max, ok := zipEstimateAgg.Buckets[0].Aggregations.Max("maxfinal")

			if ok {
				zipEstimate.Max = max.Value
			}

			min, ok := zipEstimateAgg.Buckets[0].Aggregations.Min("minfinal")

			if ok {
				zipEstimate.Min = min.Value
			}
		}

		estimates[ZIP] = zipEstimate


		x, _ := b.Aggregations.GeoDistance("geographic-range-estimates")

		for _, k := range x.Buckets{
			max, _ := k.Aggregations.Max("maxprice")
			min, _ := k.Aggregations.Min("minprice")
			if k.Key == "city" {
				estimates[FITYMILES] = &ProductEstimate{Etype:FITYMILES, Min: min.Value, Max: max.Value}
			}
			if k.Key == "metro" {
				estimates[HUNDREDMILES] = &ProductEstimate{Etype:HUNDREDMILES, Min: min.Value, Max: max.Value}
			}
			if k.Key == "national" {
				estimates[NATIONALMILES] = &ProductEstimate{Etype:NATIONALMILES, Min: min.Value, Max: max.Value}
			}
		}

		responses[b.KeyNumber] = estimates

	}

	return responses

}

func (svc *productService) GetProducts(filter ProductFilter) ([]Product, *errors.ApiError) {
	ctx := context.Background()
	filters := []elastic.Query{}
	fq := elastic.NewFunctionScoreQuery()

	priceStatusFilter := elastic.NewTermsQuery("priceStatus", "0")
	productStatusFilter := elastic.NewTermsQuery("productStatus", "0")

	keywordMultiMatch := elastic.NewMultiMatchQuery(filter.keyword, "productDesc.raw^5", "brandDesc.raw^3", "typeDesc.raw^10", "fullProductName^10", "tags")

	filters = append(filters, productStatusFilter)
	filters = append(filters, priceStatusFilter)

	if filter.location != nil {
		geoLocationFilter := elastic.NewGeoDistanceQuery("location").GeoPoint(filter.location).Distance("25mi")
		filters = append(filters, geoLocationFilter)
	}

	if filter.productIds != nil {
		productIdFilter := elastic.NewTermsQuery("productId", filter.productIds)
		filters = append(filters, productIdFilter)
	}

	if filter.categoryId != nil {
		categoryIdFilter := elastic.NewTermsQuery("categoryId", *filter.categoryId)
		filters = append(filters, categoryIdFilter)
	}

	if filter.brandId != nil {
		brandIdFilter := elastic.NewTermsQuery("brandId", *filter.brandId)
		filters = append(filters, brandIdFilter)
	}

	if filter.typeId != nil {
		typeIdFilter := elastic.NewTermsQuery("typeId", *filter.typeId)
		filters = append(filters, typeIdFilter)
	}

	businessValueScore := elastic.NewFieldValueFactorFunction().Field("typeBusinessValue").Missing(1.0)
	popularityScore := elastic.NewFieldValueFactorFunction().Field("popularity").Factor(0.001).Missing(1.0)
	totalPricesScore := elastic.NewFieldValueFactorFunction().Field("totalPricesScore").Factor(0.0001).Missing(1.0)

	bq := elastic.NewBoolQuery().Should(keywordMultiMatch).Filter(filters...)

	s := []string{
		"productId",
		"storeId",
		"listPricePerOne",
		"listPrice",
		"listQuantity",
		"listPriceUserId",
		"salePricePerOne",
		"salePriceTypeId",
		"saleValue1",
		"saleValue2",
		"saleValue3",
		"salePriceUserId",
		"saleEndDate",
		"couponPricePerOne",
		"couponPriceTypeId",
		"couponValue1",
		"couponValue2",
		"couponValue3",
		"couponUri",
		"couponEndDate",
		"unitPricePerOne",
		"unitDesc",
		"finalPrice",
	}

	finalPriceInnerHit := elastic.NewInnerHit().
		Name("MINIMAL").
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("productId")).
		DocvalueFields(s...).
		Sort("finalPrice", true).
		Size(1)

	listPricePerOneInnerHit := elastic.NewInnerHit().
		Name("MAXIMUM").
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("productId")).
		DocvalueFields(s...).
		Sort("listPricePerOne", false).
		Size(1)

	collapseBuilder := elastic.
		NewCollapseBuilder("productId").
		InnerHit(finalPriceInnerHit).
		InnerHit(listPricePerOneInnerHit)

	searchResult, searchErr := config.ElasticClient.
		Search("prices").
		Query(fq.
			Query(bq).
			AddScoreFunc(businessValueScore).
			AddScoreFunc(popularityScore).
			AddScoreFunc(totalPricesScore).
			ScoreMode("multiply").
			BoostMode("multiply")).
		Collapse(collapseBuilder).
		Do(ctx)

	if searchErr != nil {
		return nil, errors.ServerError(searchErr)
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

func (svc *productService) GetProductPrices(filter ProductPriceFilter) ([]Product, *errors.ApiError) {

	ctx := context.Background()

	fq := elastic.NewFunctionScoreQuery()

	productIdQuery := elastic.NewTermQuery("productId", filter.productId)
	priceStatusFilter := elastic.NewTermsQuery("priceStatus", "0")
	productStatusFilter := elastic.NewTermsQuery("productStatus", "0")

	bq := elastic.NewBoolQuery().Filter(productIdQuery).Filter(priceStatusFilter, productStatusFilter)

	if len(filter.storeIds) > 0 {
		s := make([]interface{}, len(filter.storeIds))
		for i, v := range filter.storeIds {
			s[i] = v
		}

		storeIdFilter := elastic.NewTermsQuery("storeId", s...)
		bq.Filter(storeIdFilter)
	}

	businessValueScore := elastic.NewFieldValueFactorFunction().Field("typeBusinessValue").Missing(1.0)
	popularityScore := elastic.NewFieldValueFactorFunction().Field("popularity").Factor(0.001).Missing(1.0)
	totalPricesScore := elastic.NewFieldValueFactorFunction().Field("totalPricesScore").Factor(0.0001).Missing(1.0)

	if filter.location != nil {
		geoLocationFilter := elastic.NewGeoDistanceQuery("location").GeoPoint(filter.location).Distance("25mi")
		bq.Filter(geoLocationFilter)
	}

	fq.Query(bq).AddScoreFunc(businessValueScore).AddScoreFunc(popularityScore).AddScoreFunc(totalPricesScore)

	searchResult, err := config.ElasticClient.Search("prices").Query(fq).Do(ctx)

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
