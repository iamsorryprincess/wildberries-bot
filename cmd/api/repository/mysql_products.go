package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type MysqlProductRepository struct {
	logger log.Logger
	conn   *mysql.Connection
}

func NewMysqlProductRepository(logger log.Logger, conn *mysql.Connection) *MysqlProductRepository {
	return &MysqlProductRepository{
		logger: logger,
		conn:   conn,
	}
}

func (r *MysqlProductRepository) Update(ctx context.Context, products []model.Product) error {
	const insertProductsSQL = `insert into
  products (
    id,
	category_id,
    name,
    rating,
    url,
    brand,
    brand_id,
    colors,
    created_at
  ) values `
	const insertProductsValuesStmt = "(?, ?, ?, ?, ?, ?, ?, ?, NOW())"

	const insertProductSizesSQL = `insert into
  products_sizes (
    product_id,
    name,
    first_price,
    previous_price,
    current_price,
    created_at
  ) values `
	const insertProductSizesValuesStmt = "(?, ?, ?, ?, ?, NOW())"

	var insertProductsBuilder strings.Builder
	insertProductsBuilder.WriteString(insertProductsSQL)
	productArgs := make([]interface{}, 0, len(products)*8)

	var insertSizesBuilder strings.Builder
	insertSizesBuilder.WriteString(insertProductSizesSQL)
	var sizeArgs []interface{}

	sizesIndex := 0
	for i, product := range products {
		if i > 0 {
			insertProductsBuilder.WriteString(", ")
		}

		colorsJSON, err := json.Marshal(product.Colors)
		if err != nil {
			r.logger.Warn().Err(err).Uint64("product_id", product.ID).Msg("failed marshal colors array")
			colorsJSON = []byte("[]")
		}

		insertProductsBuilder.WriteString(insertProductsValuesStmt)
		productArgs = append(productArgs, product.ID, product.CategoryID, product.Name, product.Rating, product.URL, product.Brand, product.BrandID, string(colorsJSON))

		for _, size := range product.Sizes {
			if sizesIndex > 0 {
				insertSizesBuilder.WriteString(", ")
			}

			insertSizesBuilder.WriteString(insertProductSizesValuesStmt)
			sizeArgs = append(sizeArgs, product.ID, size.Name, size.CurrentPrice, size.CurrentPrice, size.CurrentPrice)
			sizesIndex++
		}
	}

	const duplicateProductsStmt = " on duplicate key update updated_at = NOW();"
	insertProductsBuilder.WriteString(duplicateProductsStmt)

	productsQuery := insertProductsBuilder.String()
	if _, err := r.conn.ExecContext(ctx, productsQuery, productArgs...); err != nil {
		return fmt.Errorf("mysql products repository: failed exec insert products: %w", err)
	}

	const duplicateSizesStmt = ` as new_values on duplicate key update
  previous_price = products_sizes.current_price,
  current_price = new_values.current_price,
  updated_at = NOW();`
	insertSizesBuilder.WriteString(duplicateSizesStmt)

	sizesQuery := insertSizesBuilder.String()
	if _, err := r.conn.ExecContext(ctx, sizesQuery, sizeArgs...); err != nil {
		return fmt.Errorf("mysql products repository: failed exec insert product sizes: %w", err)
	}

	return nil
}

func (r *MysqlProductRepository) GetSizes(ctx context.Context, category string) ([]string, error) {
	const query = `select ps.name, count(ps.name) as c
from products_sizes as ps
left join products as p on p.id = ps.product_id
left join categories as cat on cat.id = p.category_id
where cat.name = ?
group by ps.name
having c >= ?;`

	const itemsCount = 100

	rows, err := r.conn.QueryContext(ctx, query, category, itemsCount)
	if err != nil {
		return nil, fmt.Errorf("mysql products repository: failed get sizes: %w", err)
	}

	defer r.conn.CloseRows(rows)

	var sizes []string
	for rows.Next() {
		var size string
		var count int

		if err = rows.Scan(&size, &count); err != nil {
			return nil, fmt.Errorf("mysql products repository: failed scan get sizes row: %w", err)
		}

		sizes = append(sizes, size)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("mysql products repository: get sizes rows error: %w", err)
	}

	return sizes, nil
}
