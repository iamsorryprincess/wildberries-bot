package mysql

import (
	"database/sql"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/background"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type Connection struct {
	logger log.Logger
	config Config
	*sql.DB
}

func NewConnection(logger log.Logger, config Config, closerStack background.CloserStack) (*Connection, error) {
	db, err := sql.Open("mysql", config.ConnectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetConnMaxLifetime(config.ConnectionMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	connection := &Connection{
		logger: logger,
		config: config,
		DB:     db,
	}

	closerStack.Push(connection)

	return connection, nil
}

func (c *Connection) CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		c.logger.Error().Err(err).Msg("mysql failed close rows")
	}
}
