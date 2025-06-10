package products

type ProductItem struct {
	Name     string   `json:"name"`
	Quantity []string `json:"quantity"`
}

type ProductKapibara struct {
	Cook    map[string][]ProductItem `json:"cook"`
	Cashier map[string][]ProductItem `json:"cashier"`
}
type ProductFresfcoff struct {
	Products map[string][]ProductItem `json:"products"`
}
type ProductsConfig struct {
	Kapibara  *ProductKapibara
	Fresfcoff *ProductFresfcoff
}
