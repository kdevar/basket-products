package stores

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kdevar/basket-products/errors"
	"github.com/kdevar/basket-products/util"
	"github.com/olivere/elastic"
	"github.com/kdevar/basket-products/const"
)

type storesServiceImpl struct {
	elasticClient *elastic.Client
}

func GetLatLonString(r *elastic.GeoPoint) string {
	return util.ConvertPointToString(r)
}

func (svc *storesServiceImpl) GetStoresForLocation(point *elastic.GeoPoint) ([]Store, *errors.ApiError) {

	origin := GetLatLonString(point)

	priceStatusFilter := elastic.NewTermsQuery(_const.PRICESTATUSFIELD, "0")
	productStatusFilter := elastic.NewTermsQuery(_const.PRODUCTSTATUSFIELD, "0")

	geoQuery := elastic.
		NewGeoDistanceQuery("location").
		GeoPoint(point).
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

	result, err := svc.elasticClient.
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
		for _, hit := range ts.Hits.Hits {
			json.Unmarshal(*hit.Source, &store)
			stores = append(stores, store)
		}
	}
	return stores, nil
}

func (svc *storesServiceImpl) GetChainsForLocation(point *elastic.GeoPoint) ([]Chain, *errors.ApiError) {
	stores, err := svc.GetStoresForLocation(point)

	if err != nil {
		return nil, err
	}

	chains := make(map[int]Chain)
	chainsArr := []Chain{}

	for _, store := range stores {
		if val, ok := chains[store.ChainID]; !ok {
			chains[store.ChainID] = store.Chain
			chainsArr = append(chainsArr, val)
		}
	}

	return chainsArr, nil
}
