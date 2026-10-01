package usecase

import (
	"context"

	"github.com/rom6621/microservices-demo/services/productcatalogservice/internal/domain"
)

type ProductRepository interface {
	FindAll(context.Context) ([]*domain.Product, error)
	FindByQuery(context.Context, string) ([]*domain.Product, error)
	FindByID(context.Context, string) (*domain.Product, error)
}

type ProductUsecase struct {
	pr ProductRepository
}

func NewProductUsecase(pr ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		pr: pr,
	}
}

func (u *ProductUsecase) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	ps, err := u.pr.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (u *ProductUsecase) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	p, err := u.pr.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (u *ProductUsecase) SearchProducts(ctx context.Context, query string) ([]*domain.Product, error) {
	ps, err := u.pr.FindByQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	return ps, nil
}
