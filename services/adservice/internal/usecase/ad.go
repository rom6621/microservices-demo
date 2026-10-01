package usecase

import (
	"context"

	"github.com/rom6621/microservice-demo/services/adservice/internal/domain"
)

type AdRepository interface {
	FindByRandom(context.Context) ([]*domain.Ad, error)
	FindByCategory(context.Context, string) ([]*domain.Ad, error)
}

type AdUsecase struct {
	ar AdRepository
}

func NewAdUsecase(ar AdRepository) *AdUsecase {
	return &AdUsecase{
		ar: ar,
	}
}

func (u *AdUsecase) GetAds(ctx context.Context, contextKeys []string) ([]*domain.Ad, error) {
	allAds := []*domain.Ad{}
	for _, k := range contextKeys {
		ads, err := u.ar.FindByCategory(ctx, k)
		if err != nil {
			return nil, err
		}
		allAds = append(allAds, ads...)
	}

	if len(allAds) == 0 {
		return u.ar.FindByRandom(ctx)
	}

	return allAds, nil
}
