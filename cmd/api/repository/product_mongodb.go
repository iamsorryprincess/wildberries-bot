package repository

import (
	"context"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
)

type MongodbProductRepository struct{}

func NewMongodbProductRepository() *MongodbProductRepository {
	return &MongodbProductRepository{}
}

func (r *MongodbProductRepository) Update(_ context.Context, _ []model.Product) error {
	return nil
}
