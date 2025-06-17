package repository

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type MysqlTrackingRepository struct {
	logger log.Logger
	conn   *mysql.Connection
}

func NewMysqlTrackingRepository(logger log.Logger, conn *mysql.Connection) *MysqlTrackingRepository {
	return &MysqlTrackingRepository{
		logger: logger,
		conn:   conn,
	}
}

func (r *MysqlTrackingRepository) AddTracking(ctx context.Context, settings model.TrackingSettings) error {
	const query = `insert into
  tracking_settings (chat_id, size_id, category_id, diff_value)
values
  (?, ?, ?, ?) as new_values on duplicate key
update
  diff_value = new_values.diff_value,
  updated_at = NOW();`

	_, err := r.conn.ExecContext(ctx, query, settings.ChatID, settings.SizeID, settings.CategoryID, settings.DiffValue)
	if err != nil {
		return fmt.Errorf("mysql insert tracking_settings error: %w", err)
	}

	return nil
}

func (r *MysqlTrackingRepository) FindMatchTracking(ctx context.Context, category string) ([]model.TrackingResult, error) {
	const query = `select
  p.name,
  p.url,
  ps.name,
  FORMAT(ps.previous_price, 2) as previous_price,
  FORMAT(ps.current_price, 2) as current_price,
  ROUND(((ps.previous_price - ps.current_price) / ps.previous_price * 100)) as diff_percent,
  ts.chat_id
from
  products_sizes as ps
  join products as p on p.id = ps.product_id
  join tracking_settings as ts on ts.size = ps.name
where
  ts.category = ?
  AND ROUND(((ps.previous_price - ps.current_price) / ps.previous_price * 100)) >= ts.diff_value;`

	rows, err := r.conn.QueryContext(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("mysql query find match tracking error: %w", err)
	}

	defer r.conn.CloseRows(rows)

	var result []model.TrackingResult
	for rows.Next() {
		var trackingResult model.TrackingResult

		if err = rows.Scan(
			&trackingResult.ProductName,
			&trackingResult.ProductURL,
			&trackingResult.Size,
			&trackingResult.PreviousPrice,
			&trackingResult.CurrentPrice,
			&trackingResult.DiffPercent,
			&trackingResult.ChatID,
		); err != nil {
			return nil, fmt.Errorf("mysql scan find match tracking row error: %w", err)
		}

		result = append(result, trackingResult)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("mysql find match tracking rows error: %w", err)
	}

	return result, nil
}
