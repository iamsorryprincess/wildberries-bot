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

func (r *MysqlTrackingRepository) FindMatchTracking(ctx context.Context, categoryID uint64) ([]model.TrackingResult, error) {
	const query = `select
  ps.product_id,
  p.name,
  p.url,
  ps.size_id,
  s.name,
  ps.previous_price,
  ps.current_price,
  ps.current_price_int,
  ROUND(((ps.previous_price - ps.current_price) / ps.previous_price * 100)) as diff_percent,
  ts.chat_id
from
  products_sizes as ps
  join products as p on p.id = ps.product_id
  join tracking_settings as ts on ts.size_id = ps.size_id
  join sizes as s on s.id = ts.size_id
  left join tracking_logs as tl on tl.chat_id = ts.chat_id and tl.size_id = ts.size_id and tl.product_id = ps.product_id
where
  ts.category_id = ? and
  (tl.price is NULL or tl.price <> ps.current_price_int) and
  ROUND(((ps.previous_price - ps.current_price) / ps.previous_price * 100)) >= ts.diff_value;`

	rows, err := r.conn.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("mysql query find match tracking error: %w", err)
	}

	defer r.conn.CloseRows(rows)

	var result []model.TrackingResult
	for rows.Next() {
		var trackingResult model.TrackingResult

		if err = rows.Scan(
			&trackingResult.ProductID,
			&trackingResult.ProductName,
			&trackingResult.ProductURL,
			&trackingResult.SizeID,
			&trackingResult.Size,
			&trackingResult.PreviousPrice,
			&trackingResult.CurrentPrice,
			&trackingResult.CurrentPriceInt,
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

func (r *MysqlTrackingRepository) SaveTrackingLog(ctx context.Context, log model.TrackingLog) error {
	const query = `insert into
  tracking_logs (chat_id, size_id, product_id, price)
values
  (?, ?, ?, ?) as new_values on duplicate key
update
  price = new_values.price,
  updated_at = NOW();`

	_, err := r.conn.ExecContext(ctx, query, log.ChatID, log.SizeID, log.ProductID, log.Price)
	if err != nil {
		return fmt.Errorf("mysql insert tracking error: %w", err)
	}

	return nil
}

func (r *MysqlTrackingRepository) DeleteTrackingSettingsByChat(ctx context.Context, chatID int64) error {
	const query = "delete from tracking_settings where chat_id = ?"
	if _, err := r.conn.ExecContext(ctx, query, chatID); err != nil {
		return fmt.Errorf("mysql delete from tracking_settings error: %w", err)
	}
	return nil
}

func (r *MysqlTrackingRepository) DeleteTrackingSettings(ctx context.Context, chatID int64, sizeID uint64, categoryID uint64) error {
	const query = "delete from tracking_settings where chat_id = ? and size_id = ? and category_id = ?"
	if _, err := r.conn.ExecContext(ctx, query, chatID, sizeID, categoryID); err != nil {
		return fmt.Errorf("mysql delete from tracking_settings error: %w", err)
	}
	return nil
}

func (r *MysqlTrackingRepository) GetTrackingSettingsInfo(ctx context.Context, chatID int64) ([]model.TrackingSettingsInfo, error) {
	const query = `select
  ts.chat_id,
  ts.size_id,
  s.name,
  ts.category_id,
  c.title,
  c.emoji,
  ts.diff_value
from
  tracking_settings as ts
  join sizes as s on s.id = ts.size_id
  join categories as c on c.id = ts.category_id
where
  chat_id = ?;`

	rows, err := r.conn.QueryContext(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("mysql get tracking settings info error: %w", err)
	}

	defer r.conn.CloseRows(rows)

	var result []model.TrackingSettingsInfo
	for rows.Next() {
		var trackingSettingsInfo model.TrackingSettingsInfo

		if err = rows.Scan(
			&trackingSettingsInfo.ChatID,
			&trackingSettingsInfo.SizeID,
			&trackingSettingsInfo.Size,
			&trackingSettingsInfo.CategoryID,
			&trackingSettingsInfo.CategoryTitle,
			&trackingSettingsInfo.CategoryEmoji,
			&trackingSettingsInfo.DiffPercent,
		); err != nil {
			return nil, fmt.Errorf("mysql scan tracking settings info row error: %w", err)
		}

		result = append(result, trackingSettingsInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("mysql get tracking settings info rows error: %w", err)
	}

	return result, nil
}
