package storage

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type DbInfo struct {
	Host  string `json:"host" yaml:"host"`
	Port  string `json:"port" yaml:"port"`
	Login string `json:"login" yaml:"login"`
	Pass  string `json:"pass" yaml:"pass"`
	Db    string `json:"db" yaml:"db"`
}

func New(dbInfo *DbInfo) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbInfo.Host, dbInfo.Port, dbInfo.Login, dbInfo.Pass, dbInfo.Db,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %v", err)
	}

	return db, nil
}
