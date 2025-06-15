package service

import (
	"context"
	"errors"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type ProductClient interface {
	GetProducts(ctx context.Context, request model.ProductsRequest) ([]model.Product, error)
}

type CategoryRepository interface {
	GetCategory(ctx context.Context, id uint64) (model.Category, error)
}

type ProductUpdateRepository interface {
	Update(ctx context.Context, products []model.Product) error
}

type ProductUpdateService struct {
	logger             log.Logger
	client             ProductClient
	categoryRepository CategoryRepository
	productRepository  ProductUpdateRepository
}

func NewUpdateService(
	logger log.Logger,
	client ProductClient,
	categoryRepository CategoryRepository,
	productRepository ProductUpdateRepository,
) *ProductUpdateService {
	return &ProductUpdateService{
		logger:             logger,
		client:             client,
		categoryRepository: categoryRepository,
		productRepository:  productRepository,
	}
}

func (s *ProductUpdateService) Update(ctx context.Context, categoryID uint64) error {
	category, err := s.categoryRepository.GetCategory(ctx, categoryID)
	if err != nil {
		return err
	}

	request := model.ProductsRequest{
		Page:       1,
		Category:   category.Name,
		CategoryID: categoryID,
	}

	for {
		s.logger.Debug().Int("page", request.Page).Msg("products request")

		products, err := s.client.GetProducts(ctx, request)
		if err != nil {
			if errors.Is(err, model.ErrRequestLimit) {
				return nil
			}
			return err
		}

		if len(products) == 0 {
			return nil
		}

		if err = s.productRepository.Update(ctx, products); err != nil {
			return err
		}

		request.Page++
	}
}
