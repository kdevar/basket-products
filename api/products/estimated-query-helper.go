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
	sub := GeneratePriceAggregation(_const.METROAREAIDFIELD, metroAreaId)
	a.Agg.SubAggregation(METROLABEL,sub)
	return a
}
func (a *EstimatedAggregationBuilder) AddCityIdSubAggregation(cityId string) *EstimatedAggregationBuilder{
	sub := GeneratePriceAggregation(_const.CITYIDFIELD, cityId)
	a.Agg.SubAggregation(CITYLABEL, sub)
	return a
}
func (a *EstimatedAggregationBuilder) AddZipIdSubAggregation(zipCodeId string) *EstimatedAggregationBuilder{
	sub := GeneratePriceAggregation(_const.ZIPIDFIELD, zipCodeId)
	a.Agg.SubAggregation(ZIPLABEL, sub)
	return a
}
func (a *EstimatedAggregationBuilder) AddGeographicRangeSubAggregation(point *elastic.GeoPoint) *EstimatedAggregationBuilder{
	sub := elastic.
		NewGeoDistanceAggregation().
		Field(_const.LOCATIONFIELD).
		Point(util.ConvertGeoPointToString(point)).
		Unit("miles").
		AddRangeWithKey(string(FIFTYMILES), 0, 50).
		AddRangeWithKey(string(HUNDREDMILES), 51, 100).
		AddRangeWithKey(string(NATIONALMILES), 101, 4000)
	a.Agg.SubAggregation(GEOGRAPHICLABEL, sub)


	s := elastic.NewMinAggregation().Field(_const.FINALPRICEFIELD)
	m := elastic.NewMaxAggregation().Field(_const.FINALPRICEFIELD)

	sub.SubAggregation(MINLABEL, s)
	sub.SubAggregation(MAXLABEL, m)

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
		Etype: CITY,
	}
	cityBucket, _ := e.Result.Aggregations.Terms(CITYLABEL)
	cityEstimate.transformFromTermsBucket(cityBucket)
	return cityEstimate
}
func (e *EstimationAggregationResult) GetMetroEstimate() *ProductEstimate{
	metroEstimate := &ProductEstimate{
		Etype: METRO,
	}
	metroBucket, _ := e.Result.Aggregations.Terms(METROLABEL)
	metroEstimate.transformFromTermsBucket(metroBucket)
	return metroEstimate
}
func (e *EstimationAggregationResult) GetZipEstimate() *ProductEstimate{
	zipEstimate := &ProductEstimate{
		Etype: ZIP,
	}
	zipBucket, _ := e.Result.Aggregations.Terms(ZIPLABEL)
	zipEstimate.transformFromTermsBucket(zipBucket)
	return zipEstimate
}

func (e *EstimationAggregationResult) GetRangeEstimate(key EstimateType) (*ProductEstimate){
	i, _ := e.Result.Aggregations.GeoDistance(GEOGRAPHICLABEL)
	for _, k := range i.Buckets {
		if k.Key == key.String() {
			max, _ := k.Aggregations.Max(MAXLABEL)
			min, _ := k.Aggregations.Min(MINLABEL)
			return &ProductEstimate{
				Etype: key,
				Min:   min.Value, Max:   max.Value}
		}
	}
	return nil
}



type EstimationAggregationResults struct{
	Current int
	Result *elastic.AggregationBucketKeyItems
}
func (e *EstimationAggregationResults) Value() *EstimationAggregationResult{
	result := &EstimationAggregationResult{
		Result: e.Result.Buckets[e.Current],
	}
	return result
}
func (e *EstimationAggregationResults) Next() bool{
	e.Current++
	 return e.Current < len(e.Result.Buckets)
}

func NewEstimationAggregationResults(r *elastic.AggregationBucketKeyItems) *EstimationAggregationResults {
	return &EstimationAggregationResults{
		Current: -1,
		Result: r,
	}
}

func GeneratePriceAggregation(field string, include string) *elastic.TermsAggregation {
	return elastic.
		NewTermsAggregation().
		Field(field).
		Include(include).
		SubAggregation(MINFINALPRICELABEL, elastic.NewMinAggregation().Field(_const.FINALPRICEFIELD)).
		SubAggregation(MAXFINALPRICELABEL, elastic.NewMaxAggregation().Field(_const.FINALPRICEFIELD)).
		SubAggregation(MINLISTLABEL, elastic.NewMinAggregation().Field(_const.LISTPRICEFIELD)).
		SubAggregation(MAXLISTLABEL, elastic.NewMaxAggregation().Field(_const.LISTPRICEFIELD))
}


type EstimateType string

func(e *EstimateType) String() string{
	return string(*e)
}

const (
	ZIP           EstimateType = "ZIP"
	CITY          EstimateType = "CITY"
	METRO         EstimateType = "METRO"
	FIFTYMILES    EstimateType = "FIFTYMILE"
	HUNDREDMILES  EstimateType = "HUNDREDMILES"
	NATIONALMILES EstimateType = "NATIONALMILES"
)


const (
	ZIPLABEL           string = "zipcode-estimate"
	CITYLABEL          string = "city-estimate"
	METROLABEL         string = "metro-estimate"
	GEOGRAPHICLABEL    string = "geographic-range-estimates"
	CHAINLABEL         string = "chain-groups"
	MAXLABEL           string = "maxprice"
	MINLABEL           string = "minprice"
	MAXFINALPRICELABEL string = "maxfinal"
	MINFINALPRICELABEL string = "minfinal"
	MAXLISTLABEL       string = "maxlist"
	MINLISTLABEL       string = "minlist"
)