package products

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kdevar/basket-products/api/errors"
	"github.com/kdevar/basket-products/config"
	"github.com/olivere/elastic"
	"strconv"
	"github.com/kdevar/basket-products/api/area"
)

var Service *productService

func init() {
	Service = &productService{}
}

type ProductFilter struct {
	keyword    string
	location   *elastic.GeoPoint
	productIds []int
	categoryId *string
}


type productService struct{}

type ProductEstimateRequest struct {
	ProductId string
	ChainId   string
	Location  elastic.GeoPoint
}



func (r *ProductEstimateRequest) GetLatLonString() string {
	return strconv.FormatFloat(r.Location.Lat, 'f', 6, 64) + "," + strconv.FormatFloat(r.Location.Lon, 'f', 6, 64)
}

func (svc *productService) GetProductEstimate(request ProductEstimateRequest) map[string]map[string]*json.RawMessage {
	ctx := context.Background()
	area := area.Service.GetAreaInformation(request.Location)

	productIdTermQuery := elastic.NewTermQuery("productId", request.ProductId)
	chainIdTermQuery := elastic.NewTermQuery("chainId", request.ChainId)

	gq := elastic.NewBoolQuery().Must(productIdTermQuery, chainIdTermQuery)

	ma := elastic.
		NewTermsAggregation().
		Field("metroAreaId").
		Include(strconv.Itoa(area.MetroAreaID)).
		SubAggregation("min-price", elastic.NewMinAggregation().Field("finalPrice")).
		SubAggregation("max-price", elastic.NewMaxAggregation().Field("finalPrice"))

	c := elastic.
		NewTermsAggregation().
		Field("cityId").
		Include(strconv.Itoa(area.CityID)).
		SubAggregation("min-price", elastic.NewMinAggregation().Field("finalPrice")).
		SubAggregation("max-price", elastic.NewMaxAggregation().Field("finalPrice"))

	z := elastic.
		NewTermsAggregation().
		Field("zipCodeId").
		Include(strconv.Itoa(area.PostalCodeID)).
		SubAggregation("min-price", elastic.NewMinAggregation().Field("finalPrice")).
		SubAggregation("max-price", elastic.NewMaxAggregation().Field("finalPrice"))

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

	a.SubAggregation("min-price", s)
	a.SubAggregation("max-price", m)

	searchQuery := config.
		ElasticClient.
		Search("prices").
		Query(gq).
		Aggregation("geographic-range-estimates", a).
		Aggregation("metro-estimate", ma).
		Aggregation("city-estimate", c).
		Aggregation("zipcode-estimate", z)

	searchResult, err := searchQuery.Do(ctx)

	if err != nil {
		fmt.Println(err)
	}

	cityAgg, _ := searchResult.
		Aggregations.
		Terms("city-estimate")

	metroAgg, _ := searchResult.
		Aggregations.
		Terms("metro-estimate")

	zipAgg, _ := searchResult.
		Aggregations.
		Terms("zipcode-estimate")

	fmt.Println(cityAgg)
	fmt.Println(metroAgg)

	agg, _ := searchResult.
		Aggregations.
		GeoDistance("geographic-range-estimates")


	responses := map[string]map[string]*json.RawMessage{
		"geographic-range-estimates": agg.Aggregations,
		"cityid-estimate": cityAgg.Aggregations,
		"metroareaid-estimate-":metroAgg.Aggregations,
		"zipcodeid-estimate": zipAgg.Aggregations,
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
