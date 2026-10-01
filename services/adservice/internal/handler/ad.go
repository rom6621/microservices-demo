package handler

import (
	"context"

	pb "github.com/rom6621/microservice-demo/services/adservice/genproto"
	"github.com/rom6621/microservice-demo/services/adservice/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdUsecase interface {
	GetAds(context.Context, []string) ([]*domain.Ad, error)
}

type AdHandler struct {
	pb.UnimplementedAdServiceServer
	au AdUsecase
}

func NewAdHandler(au AdUsecase) *AdHandler {
	return &AdHandler{au: au}
}

func (h *AdHandler) GetAds(ctx context.Context, req *pb.AdRequest) (*pb.AdResponse, error) {
	ads, err := h.au.GetAds(ctx, req.ContextKeys)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := make([]*pb.Ad, 0, len(ads))
	for _, ad := range ads {
		resp = append(resp, h.parseAd(ad))
	}

	return &pb.AdResponse{Ads: resp}, nil
}

func (h *AdHandler) parseAd(ad *domain.Ad) *pb.Ad {
	return &pb.Ad{
		Text:        ad.Text(),
		RedirectUrl: ad.RedirectURL().String(),
	}
}
