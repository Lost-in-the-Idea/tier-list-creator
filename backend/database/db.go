package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
	SQLDB *sql.DB
}

func (db *Database) InitialiseDatabase(DBName string, DBUser string, DBPassword string, DBHost string, DBPort string) error {
	dsn := "host=" + DBHost + " user=" + DBUser + " password=" + DBPassword + " dbname=" + DBName + " port=" + DBPort + " sslmode=disable"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db.DB = gormDB

	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	db.SQLDB = sqlDB

	err = sqlDB.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Close() error {
	return db.SQLDB.Close()
}