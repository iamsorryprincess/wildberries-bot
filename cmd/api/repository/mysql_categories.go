package repository

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type MysqlCategoryRepository struct {
	logger log.Logger
	conn   *mysql.Connection
}

func NewMysqlCategoryRepository(logger log.Logger, conn *mysql.Connection) *MysqlCategoryRepository {
	return &MysqlCategoryRepository{
		logger: logger,
		conn:   conn,
	}
}

func (r *MysqlCategoryRepository) GetCategories(ctx context.Context) ([]model.Category, error) {
	const query = "select id, name, title, emoji, request_url, product_url from categories;"

	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("mysql get categories error: %w", err)
	}

	defer r.conn.CloseRows(rows)

	var categories []model.Category
	for rows.Next() {
		var category model.Category

		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Title,
			&category.Emoji,
			&category.RequestURL,
			&category.ProductURL,
		); err != nil {
			return nil, fmt.Errorf("mysql scan categories row error: %w", err)
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("mysql categories row error: %w", err)
	}

	return categories, nil
}

func (r *MysqlCategoryRepository) GetCategory(ctx context.Context, id uint64) (model.Category, error) {
	const query = "select id, name, title, emoji, request_url, product_url from categories where id = ?;"

	var category model.Category
	if err := r.conn.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Title,
		&category.Emoji,
		&category.RequestURL,
		&category.ProductURL,
	); err != nil {
		return category, fmt.Errorf("mysql get category error: %w", err)
	}

	return category, nil
}
