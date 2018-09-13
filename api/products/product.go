package products

import (
	"github.com/olivere/elastic"
	"time"
)

type Product struct {
	TypeDesc               string           `json:"typeDesc"`
	ListSource             string           `json:"listSource"`
	StoreStatus            int              `json:"storeStatus"`
	CouponValue2           float64          `json:"couponValue2"`
	ListPricePerOne        float64          `json:"listPricePerOne"`
	CategoryDesc           string           `json:"categoryDesc"`
	SaleValue1             float64          `json:"saleValue1"`
	CouponPricePerOne      float64          `json:"couponPricePerOne"`
	CategoryID             int              `json:"categoryId"`
	SubCategoryID          int              `json:"subCategoryId"`
	CouponEffectiveDate    string           `json:"couponEffectiveDate"`
	WholesalePricePerOne   float64          `json:"wholesalePricePerOne"`
	SubCategoryDesc        string           `json:"subCategoryDesc"`
	SaleValue3             float64          `json:"saleValue3"`
	FinalPrice             float64          `json:"finalPrice"`
	CityID                 int              `json:"cityId"`
	UnitCount              int              `json:"unitCount"`
	Size                   float64          `json:"size"`
	BrandDesc              string           `json:"brandDesc"`
	StoreID                int              `json:"storeId"`
	StoreType              string           `json:"storeType"`
	CouponValue1           float64          `json:"couponValue1"`
	CouponStartDate        string           `json:"couponStartDate"`
	Location               elastic.GeoPoint `json:"location"`
	ProductDesc            string           `json:"productDesc"`
	StateID                int              `json:"stateId"`
	ListQuantity           float32          `json:"listQuantity"`
	BrandID                int              `json:"brandId"`
	TypeBusinessValue      int              `json:"typeBusinessValue"`
	TypeID                 int              `json:"typeId"`
	Barcodes               string           `json:"barcodes"`
	ManufacturerDesc       string           `json:"manufacturerDesc"`
	ListPriceUserID        int              `json:"listPriceUserId"`
	CountryID              int              `json:"countryId"`
	Tags                   []string         `json:"tags"`
	MetroAreaID            int              `json:"metroAreaId"`
	SaleValue2             float64          `json:"saleValue2"`
	SaleSource             string           `json:"saleSource"`
	UnitPricePerOne        float64          `json:"unitPricePerOne"`
	UnitDesc               string           `json:"unitDesc"`
	ListPrice              float64          `json:"listPrice"`
	CouponSource           string           `json:"couponSource"`
	SizeDesc               string           `json:"sizeDesc"`
	BrandBusinessValue     int              `json:"brandBusinessValue"`
	AvailabilityCount      string           `json:"availabilityCount"`
	SaleEffectiveDate      string           `json:"saleEffectiveDate"`
	Version                string           `json:"@version"`
	ProductStatus          int              `json:"productStatus"`
	UnitID                 int              `json:"unitId"`
	Geohash                string           `json:"geohash"`
	SaleStartDate          string           `json:"saleStartDate"`
	CouponURI              string           `json:"couponUri"`
	ManufacturerID         int              `json:"manufacturerId"`
	UnitBaseCount          float64          `json:"unitBaseCount"`
	SalePriceTypeID        int              `json:"salePriceTypeId"`
	TotalUserPrices        int              `json:"totalUserPrices"`
	SalePricePerOne        float64          `json:"salePricePerOne"`
	ChainID                int              `json:"chainId"`
	FullStoreName          string           `json:"fullStoreName"`
	SalePriceDesc          string           `json:"salePriceDesc"`
	AvailabilityUpdateDate string           `json:"availabilityUpdateDate"`
	PercentageDiscount     float64          `json:"percentageDiscount"`
	ProductID              int              `json:"productId"`
	AvailabilityStatus     int              `json:"availabilityStatus"`
	SaleEndDate            string           `json:"saleEndDate"`
	CouponPriceTypeID      int              `json:"couponPriceTypeId"`
	CouponUserID           int              `json:"couponUserId"`
	SaleUserID             int              `json:"saleUserId"`
	CouponEndDate          string           `json:"couponEndDate"`
	ZipCodeID              int              `json:"zipCodeId"`
	PriceStatus            int              `json:"priceStatus"`
	ProductBusinessValue   int              `json:"productBusinessValue"`
	CouponValue3           float64          `json:"couponValue3"`
	ChainDesc              string           `json:"chainDesc"`
	ChainCategoryID        int              `json:"chainCategoryId"`
	ListPriceEffectiveDate string           `json:"listPriceEffectiveDate"`
	CouponPriceDesc        string           `json:"couponPriceDesc"`
	ContainerID            int              `json:"containerId"`
	Timestamp              time.Time        `json:"@timestamp"`
	ContainerDesc          string           `json:"containerDesc"`
}

func (p *Product) GetSaleEndDate() (time.Time, error) {
	return stringToTime(p.SaleEndDate)
}

func (p *Product) GetSaleEffectiveDate() (time.Time, error) {
	return stringToTime(p.SaleEffectiveDate)
}

func (p *Product) GetListPriceEffectiveDate() (time.Time, error) {
	return stringToTime(p.ListPriceEffectiveDate)
}

func stringToTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", s)

}
