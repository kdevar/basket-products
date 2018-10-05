package stores

type Store struct {
	Chain
	StoreID       int
	FullStoreName string
}

type Chain struct {
	ChainID   int
	ChainDesc string
}
