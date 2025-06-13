package repository

import (
	"context"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/database/mysql"
)

type MysqlProductRepository struct {
	conn *mysql.Connection
}

func NewMysqlProductRepository(conn *mysql.Connection) *MysqlProductRepository {
	return &MysqlProductRepository{
		conn: conn,
	}
}

func (r *MysqlProductRepository) Update(ctx context.Context, products []model.Product) error {
	return nil
}
