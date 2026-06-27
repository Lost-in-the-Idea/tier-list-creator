package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB    *gorm.DB
	sqlDB *sql.DB
}

func NewDatabase(name, user, password, host, port string) (*Database, error) {
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + name + " port=" + port + " sslmode=disable"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &Database{DB: gormDB, sqlDB: sqlDB}, nil
}

func (db *Database) Close() error {
	return db.sqlDB.Close()
}