package repository

import (
	"context"
	"math/rand/v2"
	"net/url"

	"github.com/rom6621/microservice-demo/services/adservice/internal/domain"
)

const maxAdsToServe = 2

type AdRepository struct {
	adsByCategory map[string][]*domain.Ad
	allAds        []*domain.Ad
}

func NewAdRepository() (*AdRepository, error) {
	items := []struct {
		category    string
		text        string
		redirectURL string
	}{
		{"clothing", "Tank top for sale. 20% off.", "/product/66VCHSJNUP"},
		{"accessories", "Watch for sale. Buy one, get second kit for free", "/product/1YMWWN1N4O"},
		{"footwear", "Loafers for sale. Buy one, get second one for free", "/product/L9ECAV7KIM"},
		{"hair", "Hairdryer for sale. 50% off.", "/product/2ZYFJ3GM2N"},
		{"decor", "Candle holder for sale. 30% off.", "/product/0PUK6V6EV0"},
		{"kitchen", "Bamboo glass jar for sale. 10% off.", "/product/9SIQT8TOJO"},
		{"kitchen", "Mug for sale. Buy two, get third one for free", "/product/6E92ZMYYFZ"},
	}

	adsByCategory := make(map[string][]*domain.Ad)
	allAds := make([]*domain.Ad, len(items))

	for i, item := range items {
		parsed, err := url.Parse(item.redirectURL)
		if err != nil {
			return nil, err
		}
		ad, err := domain.NewAd(item.text, parsed)
		if err != nil {
			return nil, err
		}

		adsByCategory[item.category] = append(adsByCategory[item.category], ad)
		allAds[i] = ad
	}

	return &AdRepository{
		adsByCategory: adsByCategory,
		allAds:        allAds,
	}, nil
}

func (r *AdRepository) FindByCategory(ctx context.Context, category string) ([]*domain.Ad, error) {
	return r.adsByCategory[category], nil
}

func (r *AdRepository) FindByRandom(ctx context.Context) ([]*domain.Ad, error) {
	ads := make([]*domain.Ad, 0, maxAdsToServe)
	for range maxAdsToServe {
		ads = append(ads, r.allAds[rand.N(len(r.allAds))])
	}

	return ads, nil
}
