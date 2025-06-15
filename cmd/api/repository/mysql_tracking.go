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
	fmt.Println(settings)
	return nil
}
