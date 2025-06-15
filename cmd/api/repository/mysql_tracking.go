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
  tracking_settings (chat_id, size, category, diff_value)
values
  (?, ?, ?, ?) as new_values on duplicate key
update
  diff_value = new_values.diff_value,
  updated_at = NOW();`

	_, err := r.conn.ExecContext(ctx, query, settings.ChatID, settings.Size, settings.Category, settings.DiffValue)
	if err != nil {
		return fmt.Errorf("mysql insert tracking_settings error: %w", err)
	}

	return nil
}
