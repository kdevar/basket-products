package stores

import (
	"github.com/olivere/elastic"
	"github.com/kdevar/basket-products/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kdevar/basket-products/api/errors"
	"strconv"
)

var Service *storesService

func init(){
	Service = &storesService{}
}

type storesService struct {}

func GetLatLonString(r elastic.GeoPoint) string {
	return strconv.FormatFloat(r.Lat, 'f', 6, 64) + "," + strconv.FormatFloat(r.Lon, 'f', 6, 64)
}

func (svc *storesService) GetStoresForLocation(point elastic.GeoPoint) ([]Store, *errors.ApiError){

	origin := GetLatLonString(point)

	priceStatusFilter := elastic.NewTermsQuery("priceStatus", "0")
	productStatusFilter := elastic.NewTermsQuery("productStatus", "0")

	geoQuery := elastic.
		NewGeoDistanceQuery("location").
		GeoPoint(&point).
		Distance("25mi")

	boolQuery := elastic.NewBoolQuery().Must(geoQuery, productStatusFilter, priceStatusFilter)

	expDecay := elastic.
		NewExponentialDecayFunction().
		FieldName("location").
		Scale("3mi").
		Origin(origin)

	funScoreQuery := elastic.
		NewFunctionScoreQuery().
		Query(boolQuery).
		BoostMode("REPLACE").
		AddScoreFunc(expDecay)

	a := elastic.
		NewTermsAggregation().
		Field("chainId").
		Order("max-score", false).
		Size(25).
		SubAggregation("max-score", elastic.NewMaxAggregation().Script(elastic.NewScript("_score"))).
		SubAggregation("top-hits", elastic.NewTopHitsAggregation().Size(1).SortBy(
			elastic.NewGeoDistanceSort("location").Point(point.Lat, point.Lon).Asc()))


	result, err :=config.
		ElasticClient.
		Search("prices").
		Query(funScoreQuery).
		Aggregation("chains", a).
		Do(context.Background())

	if err != nil {
		fmt.Println(err)
	}


	groups, _ := result.Aggregations.Terms("chains")
	var store Store
	var stores []Store
	for _, b := range groups.Buckets {
		ts, _ := b.Aggregations.TopHits("top-hits")
		for _,hit := range ts.Hits.Hits{
			json.Unmarshal(*hit.Source, &store)
			stores = append(stores, store)
		}
	}

	return stores, nil

}

func (svc *storesService) GetChainsForLocation(point elastic.GeoPoint) ([]Chain, *errors.ApiError){
	stores, err := svc.GetStoresForLocation(point)

	if err != nil {
		return nil, err
	}

	chains := make(map[int]Chain)
	chainsArr := []Chain{}

	for _,store := range stores {
		if val, ok := chains[store.ChainID]; !ok {
			chains[store.ChainID] = store.Chain
			chainsArr = append(chainsArr, val)
		}
	}

	return chainsArr, nil
}

func (svc *storesService) GetUserFavorites(){

}
