package repository

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/database/mysql"
)

type MysqlSizeRepository struct {
	conn *mysql.Connection
}

func NewMysqlSizeRepository(conn *mysql.Connection) *MysqlSizeRepository {
	return &MysqlSizeRepository{
		conn: conn,
	}
}

func (r *MysqlSizeRepository) GetSizesInfo(ctx context.Context, categoryID uint64) ([]model.SizeInfo, error) {
	const query = `select ps.size_id, s.name, count(ps.size_id) as c
from products_sizes as ps
left join products as p on p.id = ps.product_id
left join sizes as s on s.id = ps.size_id
where p.category_id = ?
group by ps.size_id
having c >= ?;`

	const itemsCount = 100

	rows, err := r.conn.QueryContext(ctx, query, categoryID, itemsCount)
	if err != nil {
		return nil, fmt.Errorf("mysql products repository: failed get sizes info: %w", err)
	}

	defer r.conn.CloseRows(rows)

	var result []model.SizeInfo
	for rows.Next() {
		var item model.SizeInfo

		if err = rows.Scan(&item.ID, &item.Name, &item.ProductsCount); err != nil {
			return nil, fmt.Errorf("mysql products repository: failed scan get sizes info row: %w", err)
		}

		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("mysql products repository: get sizes info rows error: %w", err)
	}

	return result, nil
}

func (r *MysqlSizeRepository) GetSizeCategoryInfo(ctx context.Context, sizeID uint64, categoryID uint64) (model.SizeCategoryInfo, error) {
	const query = `select
  (select name from sizes where id = ?) as size_name,
  (select title from categories where id = ?) as category_title,
  (select emoji from categories where id = ?) as category_emoji;`

	var result model.SizeCategoryInfo
	err := r.conn.QueryRowContext(ctx, query, sizeID, categoryID, categoryID).Scan(&result.Name, &result.CategoryTitle, &result.CategoryEmoji)
	if err != nil {
		return model.SizeCategoryInfo{}, fmt.Errorf("mysql get size category info error: %w", err)
	}

	return result, nil
}
