package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/background"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type WorkerPool interface {
	Run(ctx context.Context, name string, handler background.HandlerFunc)
}

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
	SendNotifications(ctx context.Context, category string) error
}

type ProductService struct {
	logger log.Logger

	workerPool WorkerPool
	client     ProductClient

	categoryRepository CategoryRepository
	productRepository  ProductUpdateRepository

	trackingNotifier TrackingNotifier
}

func NewProductService(
	logger log.Logger,
	workerPool WorkerPool,
	client ProductClient,
	categoryRepository CategoryRepository,
	productRepository ProductUpdateRepository,
	trackingNotifier TrackingNotifier,
) *ProductService {
	return &ProductService{
		logger:             logger,
		workerPool:         workerPool,
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

	for _, category := range categories {
		s.workerPool.Run(ctx, fmt.Sprintf("update %s products", category.Name), func(ctx context.Context) error {
			return s.UpdateProducts(ctx, category.ID)
		})
	}

	return nil
}

func (s *ProductService) UpdateProducts(ctx context.Context, categoryID uint64) error {
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

	if err = s.trackingNotifier.SendNotifications(ctx, category.Name); err != nil {
		return err
	}

	return nil
}
