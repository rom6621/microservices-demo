package handler

import (
	"context"
	"errors"

	pb "github.com/rom6621/microservices-demo/services/productcatalogservice/genproto"
	"github.com/rom6621/microservices-demo/services/productcatalogservice/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductUsecase interface {
	ListProducts(context.Context) ([]*domain.Product, error)
	SearchProducts(context.Context, string) ([]*domain.Product, error)
	GetProduct(context.Context, string) (*domain.Product, error)
}

type ProductHandler struct {
	pb.UnimplementedProductCatalogServiceServer
	pu ProductUsecase
}

func NewProductHandler(
	pu ProductUsecase,
) *ProductHandler {
	return &ProductHandler{
		pu: pu,
	}
}

func (h *ProductHandler) ListProducts(ctx context.Context, _ *pb.Empty) (*pb.ListProductsResponse, error) {
	ps, err := h.pu.ListProducts(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	products := make([]*pb.Product, 0, len(ps))
	for _, p := range ps {
		products = append(products, h.parseProduct(p))
	}

	return &pb.ListProductsResponse{Products: products}, nil
}

func (h *ProductHandler) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	ps, err := h.pu.SearchProducts(ctx, req.Query)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	products := make([]*pb.Product, 0, len(ps))
	for _, p := range ps {
		products = append(products, h.parseProduct(p))
	}

	return &pb.SearchProductsResponse{Results: products}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	p, err := h.pu.GetProduct(ctx, req.Id)
	if errors.Is(err, domain.ErrProductNotFound) {
		return nil, status.Errorf(codes.NotFound, "product %q not found", req.Id)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.parseProduct(p), nil
}

func (h *ProductHandler) parseProduct(p *domain.Product) *pb.Product {
	return &pb.Product{
		Id:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Picture:     p.Picture(),
		PriceUsd: &pb.Money{
			CurrencyCode: p.Price().CurrencyCode(),
			Units:        p.Price().Units(),
			Nanos:        p.Price().Nanos(),
		},
		Categories: p.Categories(),
	}
}
