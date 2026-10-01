package domain

import (
	"fmt"
)

type Price struct {
	currencyCode string
	units        int64
	nanos        int32
}

func NewPrice(currencyCode string, units int64, nanos int32) (*Price, error) {
	return &Price{
		currencyCode: currencyCode,
		units:        units,
		nanos:        nanos,
	}, nil
}

func (p Price) CurrencyCode() string { return p.currencyCode }
func (p Price) Units() int64         { return p.units }
func (p Price) Nanos() int32         { return p.nanos }

type Product struct {
	id          string
	name        string
	description string
	picture     string
	price       Price
	categories  []string
}

func NewProduct(id, name, description, picture string, price Price, categories []string) (*Product, error) {
	if price.currencyCode != "USD" {
		return nil, fmt.Errorf("new product: %w", ErrValidation)
	}

	return &Product{id: id, name: name, description: description, picture: picture, price: price, categories: categories}, nil
}

func (p *Product) ID() string           { return p.id }
func (p *Product) Name() string         { return p.name }
func (p *Product) Description() string  { return p.description }
func (p *Product) Picture() string      { return p.picture }
func (p *Product) Price() Price         { return p.price }
func (p *Product) Categories() []string { return p.categories }
