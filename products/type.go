package products

type productItem struct {
	Name     string   `json:"name"`
	Quantity []string `json:"quantity"`
}

type productKapibara struct {
	Cook    map[string][]productItem `json:"cook"`
	Cashier map[string][]productItem `json:"cashier"`
}
type productFresfcoff struct {
	Products map[string][]productItem `json:"products"`
}
type ProductsConfig struct {
	Kapibara  productKapibara
	Fresfcoff productFresfcoff
}
