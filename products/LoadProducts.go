package products

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadProducts(pathK, pathF string) (*ProductsConfig, error) {
	k, e := loadKapibara(pathK)
	if e != nil {
		return nil, e
	}
	f, e := loadFreshcoff(pathF)
	if e != nil {
		return nil, e
	}
	return &ProductsConfig{
		Kapibara:  k,
		Fresfcoff: f,
	}, nil
}
func loadFreshcoff(path string) (*ProductFresfcoff, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Can't read productsFreshcoff.json: ", err)
	}
	var productsFreshcoff ProductFresfcoff
	if err := json.Unmarshal(file, &productsFreshcoff); err != nil {
		return nil, fmt.Errorf("failed to unmarshal products.json: %w", err)
	}
	if len(productsFreshcoff.Products) == 0 {
		return nil, fmt.Errorf("No products found in %s", path)
	}
	return &productsFreshcoff, nil
}
func loadKapibara(path string) (*ProductKapibara, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Can't read productsKapibara.json: ", err)
	}

	var productsKapibara ProductKapibara
	if err := json.Unmarshal(file, &productsKapibara); err != nil {
		return nil, fmt.Errorf("failed to unmarshal productsKapibara.json: %w", err)
	}
	return &productsKapibara, nil
}
