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
    size_id,
    first_price,
    previous_price,
    current_price,
    created_at
  ) values `
	const insertProductSizesValuesStmt = "(?, ?, ?, ?, ?, NOW())"

	sizesMap, err := r.updateSizes(ctx, products)
	if err != nil {
		return err
	}

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
			sizeArgs = append(sizeArgs, product.ID, sizesMap[size.Name], size.CurrentPrice, size.CurrentPrice, size.CurrentPrice)
			sizesIndex++
		}
	}

	const duplicateProductsStmt = " on duplicate key update updated_at = NOW();"
	insertProductsBuilder.WriteString(duplicateProductsStmt)

	if _, err := r.conn.ExecContext(ctx, insertProductsBuilder.String(), productArgs...); err != nil {
		return fmt.Errorf("mysql products repository: failed exec insert products: %w", err)
	}

	const duplicateSizesStmt = ` as new_values on duplicate key update
  previous_price = products_sizes.current_price,
  current_price = new_values.current_price,
  updated_at = NOW();`
	insertSizesBuilder.WriteString(duplicateSizesStmt)

	if _, err := r.conn.ExecContext(ctx, insertSizesBuilder.String(), sizeArgs...); err != nil {
		return fmt.Errorf("mysql products repository: failed exec insert product sizes: %w", err)
	}

	return nil
}

func (r *MysqlProductRepository) updateSizes(ctx context.Context, products []model.Product) (map[string]uint64, error) {
	const insertQuery = "insert into sizes (name, created_at) values "
	const onDuplicateStmt = " on duplicate key update updated_at = NOW()"
	const selectQuery = "select id, name from sizes where name in ("

	var insertBuilder strings.Builder
	var selectBuilder strings.Builder

	var args []interface{}
	sizeMap := make(map[string]uint64)
	insertBuilder.WriteString(insertQuery)
	selectBuilder.WriteString(selectQuery)
	i := 0

	for _, product := range products {
		for _, size := range product.Sizes {
			_, ok := sizeMap[size.Name]
			if !ok {
				sizeMap[size.Name] = 0

				if i > 0 {
					insertBuilder.WriteString(", ")
					selectBuilder.WriteString(", ")
				}

				insertBuilder.WriteString("(?, NOW())")
				selectBuilder.WriteString("?")
				args = append(args, size.Name)
				i++
			}
		}
	}

	insertBuilder.WriteString(onDuplicateStmt)

	if len(sizeMap) > 1 {
		selectBuilder.WriteString(")")
	}

	if _, err := r.conn.ExecContext(ctx, insertBuilder.String(), args...); err != nil {
		return nil, fmt.Errorf("mysql insert sizes error: %w", err)
	}

	rows, err := r.conn.QueryContext(ctx, selectBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("mysql select sizes error: %w", err)
	}

	defer r.conn.CloseRows(rows)

	for rows.Next() {
		var id uint64
		var name string

		if err = rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("mysql scan sizes row error: %w", err)
		}

		sizeMap[name] = id
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("mysql select sizes rows error: %w", err)
	}

	return sizeMap, nil
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
