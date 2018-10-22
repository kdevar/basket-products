package products

import (
	"github.com/olivere/elastic"
	"github.com/kdevar/basket-products/const"
	"github.com/kdevar/basket-products/util"
)

type EstimationQueryBuilder struct {
	Query *elastic.BoolQuery
}
func (q *EstimationQueryBuilder) AddProductIdFilter(productId string) *EstimationQueryBuilder {
	q.Query.Filter(elastic.NewTermQuery(_const.PRODUCTIDFIELD, productId))
	return q
}
func (q *EstimationQueryBuilder) AddPriceStatusFilter() *EstimationQueryBuilder {
	q.Query.Filter(elastic.NewTermsQuery(_const.PRICESTATUSFIELD, "0"))
	return q
}
func (q *EstimationQueryBuilder) AddProductStatusFilter() *EstimationQueryBuilder {
	q.Query.Filter(elastic.NewTermsQuery(_const.PRODUCTSTATUSFIELD, "0"))
	return q
}
func (q *EstimationQueryBuilder) AddChainIdFilters(chainIds []string) *EstimationQueryBuilder {
	s := util.ConvertStringToInterface(chainIds)
	chainIdTermQuery := elastic.NewTermsQuery("chainId", s...)
	q.Query.Filter(chainIdTermQuery)
	return q
}
func (q *EstimationQueryBuilder) GetQuery() *elastic.BoolQuery{
	return q.Query
}

func NewEstimatedQuery() *EstimationQueryBuilder {
	return &EstimationQueryBuilder{
		Query: elastic.NewBoolQuery(),
	}
}

type EstimatedAggregationBuilder struct {
	Agg *elastic.TermsAggregation
}
func (a *EstimatedAggregationBuilder) AddMetroAreaSubAggregation(metroAreaId string) *EstimatedAggregationBuilder{
	sub := NewEstimateTermsAgg(_const.METROAREAIDFIELD, metroAreaId).Generate()
	a.Agg.SubAggregation(_const.METROLABEL,sub)
	return a
}
func (a *EstimatedAggregationBuilder) AddCityIdSubAggregation(cityId string) *EstimatedAggregationBuilder{
	sub := NewEstimateTermsAgg(_const.CITYIDFIELD, cityId).Generate()
	a.Agg.SubAggregation(_const.CITYLABEL, sub)
	return a
}
func (a *EstimatedAggregationBuilder) AddZipIdSubAggregation(zipCodeId string) *EstimatedAggregationBuilder{
	sub := NewEstimateTermsAgg(_const.ZIPIDFIELD, zipCodeId).Generate()
	a.Agg.SubAggregation(_const.ZIPLABEL, sub)
	return a
}
func (a *EstimatedAggregationBuilder) AddGeographicRangeSubAggregation(point *elastic.GeoPoint) *EstimatedAggregationBuilder{
	sub := elastic.
		NewGeoDistanceAggregation().
		Field(_const.LOCATIONFIELD).
		Point(util.ConvertGeoPointToString(point)).
		Unit("miles").
		AddRangeWithKey(string(_const.FIFTYMILES), 0, 50).
		AddRangeWithKey(string(_const.HUNDREDMILES), 51, 100).
		AddRangeWithKey(string(_const.NATIONALMILES), 101, 4000)
	a.Agg.SubAggregation(_const.GEOGRAPHICLABEL, sub)


	s := elastic.NewMinAggregation().Field(_const.FINALPRICEFIELD)
	m := elastic.NewMaxAggregation().Field(_const.FINALPRICEFIELD)

	sub.SubAggregation(_const.MINLABEL, s)
	sub.SubAggregation(_const.MAXLABEL, m)

	return a
}
func (a *EstimatedAggregationBuilder) GetAggregation() *elastic.TermsAggregation{
	return a.Agg
}
func NewEstimationAggregation() *EstimatedAggregationBuilder{
	return &EstimatedAggregationBuilder{
		elastic.NewTermsAggregation().Field(_const.CHAINIDFIELD),
	}
}

type EstimationAggregationResult struct{
	Result *elastic.AggregationBucketKeyItem
}
func (e *EstimationAggregationResult) GetCityEstimate() *ProductEstimate{
	cityEstimate := &ProductEstimate{
		Etype: _const.CITY,
	}
	cityBucket, _ := e.Result.Aggregations.Terms(_const.CITYLABEL)
	cityEstimate.transformFromTermsBucket(cityBucket)
	return cityEstimate
}
func (e *EstimationAggregationResult) GetMetroEstimate() *ProductEstimate{
	metroEstimate := &ProductEstimate{
		Etype: _const.METRO,
	}
	metroBucket, _ := e.Result.Aggregations.Terms(_const.METROLABEL)
	metroEstimate.transformFromTermsBucket(metroBucket)
	return metroEstimate
}
func (e *EstimationAggregationResult) GetZipEstimate() *ProductEstimate{
	zipEstimate := &ProductEstimate{
		Etype: "ZIP",
	}
	zipBucket, _ := e.Result.Aggregations.Terms(_const.ZIPLABEL)
	zipEstimate.transformFromTermsBucket(zipBucket)
	return zipEstimate
}

func (e *EstimationAggregationResult) GetRangeEstimate(key _const.EstimateType) (*ProductEstimate){
	i, _ := e.Result.Aggregations.GeoDistance(_const.GEOGRAPHICLABEL)
	for _, k := range i.Buckets {
		if k.Key == key.String() {
			max, _ := k.Aggregations.Max(_const.MAXLABEL)
			min, _ := k.Aggregations.Min(_const.MINLABEL)
			return &ProductEstimate{
				Etype: key,
				Min:   min.Value, Max:   max.Value}
		}
	}
	return nil
}



type EstimationAggregationResultIterator struct{
	Current int
	Result *elastic.AggregationBucketKeyItems
}
func (e *EstimationAggregationResultIterator) Value() *EstimationAggregationResult{
	result := &EstimationAggregationResult{
		Result: e.Result.Buckets[e.Current],
	}
	e.Current++
	return result
}
func (e *EstimationAggregationResultIterator) Next() bool{

	 return e.Current < len(e.Result.Buckets)
}

func NewEstimationAggregationResultIterator(r *elastic.AggregationBucketKeyItems) *EstimationAggregationResultIterator{
	return &EstimationAggregationResultIterator{
		Current: 0,
		Result: r,
	}
}