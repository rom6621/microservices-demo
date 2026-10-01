package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/rom6621/microservices-demo/services/productcatalogservice/internal/domain"
)

type PriceUSD struct {
	CurrencyCode string `json:"currencyCode"`
	Units        int64  `json:"units"`
	Nanos        int32  `json:"nanos"`
}
type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	PriceUSD    PriceUSD `json:"priceUsd"`
	Categories  []string `json:"categories"`
}
type CatalogJSON struct {
	Products []Product `json:"products"`
}

type ProductRepository struct {
	catalog []*domain.Product
}

func NewProductRepository(path string) (*ProductRepository, error) {
	catalogByte, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var catalogJSON CatalogJSON
	if err := json.Unmarshal(catalogByte, &catalogJSON); err != nil {
		return nil, err
	}

	catalog := make([]*domain.Product, 0, len(catalogJSON.Products))
	for _, p := range catalogJSON.Products {
		price, err := domain.NewPrice(p.PriceUSD.CurrencyCode, p.PriceUSD.Units, p.PriceUSD.Nanos)
		if err != nil {
			return nil, err
		}

		product, err := domain.NewProduct(p.ID, p.Name, p.Description, p.Picture, *price, p.Categories)
		if err != nil {
			return nil, err
		}

		catalog = append(catalog, product)
	}

	return &ProductRepository{
		catalog: catalog,
	}, nil
}

func (r *ProductRepository) FindAll(context.Context) ([]*domain.Product, error) {
	return slices.Clone(r.catalog), nil
}

func (r *ProductRepository) FindByQuery(ctx context.Context, query string) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0, len(r.catalog))
	for _, p := range r.catalog {
		if strings.Contains(strings.ToLower(p.Name()), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(p.Description()), strings.ToLower(query)) {
			products = append(products, p)
		}
	}

	return products, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	for _, p := range r.catalog {
		if p.ID() == id {
			return p, nil
		}
	}

	return nil, fmt.Errorf("find product %q: %w", id, domain.ErrProductNotFound)
}
