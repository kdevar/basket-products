package stores

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kdevar/basket-products/const"
	"github.com/kdevar/basket-products/errors"
	"github.com/kdevar/basket-products/util"
	"github.com/olivere/elastic"
)

type storesServiceImpl struct {
	elasticClient *elastic.Client
}

const (
	distance               string = "25mi"
	scale                  string = "3mi"
	boost                  string = "REPLACE"
	maxscoresubaggregation string = "max-score"
	tophitsubaggregation   string = "top-hits"
	chainsaggregation      string = "chains"
)

func (svc *storesServiceImpl) GetStoresForLocation(point *elastic.GeoPoint) ([]Store, *errors.ApiError) {

	origin := util.ConvertGeoPointToString(point)

	priceStatusFilter := elastic.NewTermsQuery(_const.PRICESTATUSFIELD, "0")
	productStatusFilter := elastic.NewTermsQuery(_const.PRODUCTSTATUSFIELD, "0")

	geoQuery := elastic.
		NewGeoDistanceQuery(_const.LOCATIONFIELD).
		GeoPoint(point).
		Distance(distance)

	boolQuery := elastic.NewBoolQuery().Must(geoQuery, productStatusFilter, priceStatusFilter)

	expDecay := elastic.
		NewExponentialDecayFunction().
		FieldName(_const.LOCATIONFIELD).
		Scale(scale).
		Origin(origin)

	funScoreQuery := elastic.
		NewFunctionScoreQuery().
		Query(boolQuery).
		BoostMode(boost).
		AddScoreFunc(expDecay)

	distanceSort := elastic.NewGeoDistanceSort(_const.LOCATIONFIELD).Point(point.Lat, point.Lon).Asc()

	chainAgg := elastic.
		NewTermsAggregation().
		Field(_const.CHAINIDFIELD).
		Order(maxscoresubaggregation, false).
		Size(25).
		SubAggregation(
			maxscoresubaggregation,
			elastic.NewMaxAggregation().Script(elastic.NewScript("_score"))).
		SubAggregation(
			tophitsubaggregation,
			elastic.NewTopHitsAggregation().Size(1).SortBy(distanceSort))

	result, err := svc.elasticClient.
		Search(_const.PROUCTPRICEINDEX).
		Query(funScoreQuery).
		Aggregation(chainsaggregation, chainAgg).
		Do(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	groups, _ := result.Aggregations.Terms(chainsaggregation)
	var store Store
	var stores []Store
	for _, b := range groups.Buckets {
		ts, _ := b.Aggregations.TopHits(tophitsubaggregation)
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
