package stores

var Service *storesService

func init(){
	Service = &storesService{}
}

type storesService struct {}

func (svc *storesService) GetStoresForLocation(){

}

func (svc *storesService) GetUserFavorites(){

}
