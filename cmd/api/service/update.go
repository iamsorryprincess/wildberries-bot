package service

import (
	"context"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type ProductClient interface {
	GetProducts(ctx context.Context, request model.ProductsRequest) ([]model.Product, error)
}

type ProductUpdateRepository interface {
	Update(ctx context.Context, products []model.Product) error
}

type ProductUpdateService struct {
	logger     log.Logger
	client     ProductClient
	repository ProductUpdateRepository
}

func NewUpdateService(logger log.Logger, client ProductClient, repository ProductUpdateRepository) *ProductUpdateService {
	return &ProductUpdateService{
		logger:     logger,
		client:     client,
		repository: repository,
	}
}

func (s *ProductUpdateService) Update(ctx context.Context) error {
	request := model.ProductsRequest{
		Page:     1,
		Category: model.ProductCategoryDresses,
	}

	for {
		s.logger.Debug().Int("page", request.Page).Msg("products request")

		products, err := s.client.GetProducts(ctx, request)
		if err != nil {
			return err
		}

		if len(products) == 0 {
			return nil
		}

		if err = s.repository.Update(ctx, products); err != nil {
			return err
		}

		request.Page++
	}
}
