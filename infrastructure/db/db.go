package db

import (
	"fmt"

	"github.com/noonbyte/platform/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	DB() *gorm.DB
	Close() error
}

type database struct {
	conn *gorm.DB
}

func New(cfg configs.DatabaseConfiguration) (Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return &database{conn: conn}, nil
}

func (d *database) DB() *gorm.DB {
	return d.conn
}

func (d *database) Close() error {
	sqlDB, err := d.conn.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	return sqlDB.Close()
}
