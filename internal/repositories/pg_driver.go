package repositories

import (
	"fmt"

	"github.com/kuzin57/auth-system/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PgDriver struct {
	db *sqlx.DB
}

func NewPgDriver(conf *config.Config) (*PgDriver, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.DBName)

	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PgDriver{db: db}, nil
}

func (d *PgDriver) DB() *sqlx.DB {
	return d.db
}
