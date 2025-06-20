package service

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type ProductClient interface {
	GetProducts(ctx context.Context, request model.ProductsRequest) ([]model.Product, error)
}

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]model.Category, error)
	GetCategory(ctx context.Context, id uint64) (model.Category, error)
}

type ProductUpdateRepository interface {
	Update(ctx context.Context, products []model.Product) error
}

type TrackingNotifier interface {
	SendNotifications(ctx context.Context, categoryID uint64) error
}

type ProductService struct {
	logger log.Logger

	hashCounter uint64

	client ProductClient

	categoryRepository CategoryRepository
	productRepository  ProductUpdateRepository

	trackingNotifier TrackingNotifier
}

func NewProductService(
	logger log.Logger,
	client ProductClient,
	categoryRepository CategoryRepository,
	productRepository ProductUpdateRepository,
	trackingNotifier TrackingNotifier,
) *ProductService {
	return &ProductService{
		logger:             logger,
		client:             client,
		categoryRepository: categoryRepository,
		productRepository:  productRepository,
		trackingNotifier:   trackingNotifier,
	}
}

func (s *ProductService) RunUpdateWorkers(ctx context.Context) error {
	categories, err := s.categoryRepository.GetCategories(ctx)
	if err != nil {
		return err
	}

	if len(categories) == 0 {
		s.logger.Warn().Msg("no product categories found")
		return nil
	}

	if len(categories) == 1 {
		return s.UpdateProducts(ctx, categories[0])
	}

	index := atomic.AddUint64(&s.hashCounter, 1) % uint64(len(categories))
	return s.UpdateProducts(ctx, categories[index])
}

func (s *ProductService) UpdateProducts(ctx context.Context, category model.Category) error {
	request := model.ProductsRequest{
		Page:       1,
		Category:   category.Name,
		CategoryID: category.ID,
		RequestURL: category.RequestURL,
		ProductURL: category.ProductURL,
	}

	var err error

	for {
		s.logger.Debug().Int("page", request.Page).Msg("products request")

		var products []model.Product
		products, err = s.client.GetProducts(ctx, request)
		if err != nil {
			break
		}

		if len(products) == 0 {
			break
		}

		if err = s.productRepository.Update(ctx, products); err != nil {
			break
		}

		request.Page++
	}

	if err != nil && !errors.Is(err, model.ErrRequestLimit) {
		return err
	}

	if err = s.trackingNotifier.SendNotifications(ctx, category.ID); err != nil {
		return err
	}

	return nil
}
