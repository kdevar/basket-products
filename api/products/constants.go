package products

type EstimateType string

const (
	ZIP           EstimateType = "ZIP"
	CITY          EstimateType = "CITY"
	METRO         EstimateType = "METRO"
	FITYMILES     EstimateType = "FIFTYMILE"
	HUNDREDMILES  EstimateType = "HUNDREDMILES"
	NATIONALMILES EstimateType = "NATIONALMILES"
)

const PROUCTPRICEINDEX string = "prices"

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

const (
	FINALPRICEFIELD        string = "finalPrice"
	BUSINESSVALUEFIELD     string = "typeBusinessValue"
	POPULARITYFIELD        string = "popularity"
	TOTALPRICESCOREFIELD   string = "totalPricesScore"
	PRODUCTIDFIELD         string = "productId"
	STOREIDFIELD           string = "storeId"
	CATEGORYIDFIELD        string = "categoryId"
	TYPEIDFIELD            string = "typeId"
	BRANDIDFIELD           string = "brandId"
	PRICESTATUSFIELD       string = "priceStatus"
	PRODUCTSTATUSFIELD     string = "productStatus"
	LOCATIONFIELD          string = "location"
	METROAREAIDFIELD       string = "metroAreaId"
	CHAINIDFIELD           string = "chainId"
	CITYIDFIELD            string = "cityId"
	ZIPIDFIELD             string = "zipCodeId"
	LISTPRICEFIELD         string = "listPrice"
	LISTPRICEPERONEFIELD   string = "listPricePerOne"
	LISTQUANTITYFIELD      string = "listQuantity"
	LISTPRICEUSERIDFIELD   string = "listPriceUserId"
	SALEPRICEPERONEFIELD   string = "salePricePerOne"
	SALEPRICETYPEID        string = "salePriceTypeId"
	SALEVALUE1FIELD        string = "saleValue1"
	SALEVALUE2FIELD        string = "saleValue2"
	SALEVALUE3FIELD        string = "saleValue3"
	SALEPRICEUSERIDFIELD   string = "salePriceUserId"
	SALEENDDATEFIELD       string = "saleEndDate"
	COUPONPRICEPERONEFIELD string = "couponPricePerOne"
	COUPONPRICETYPEIDFIELD string = "couponPriceTypeId"
	COUPONVALUE1FIELD      string = "couponValue1"
	COUPONVALUE2FIELD      string = "couponValue2"
	COUPONVALUE3FIELD      string = "couponValue3"
	COUPONURIFIELD         string = "couponUri"
	COUPONENDDATEFIELD     string = "couponEndDate"
	UNITPRICEPERONEFIELD   string = "unitPricePerOne"
	UNIDESCFIELD           string = "unitDesc"
)
