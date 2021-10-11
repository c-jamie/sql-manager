package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/c-jamie/sql-manager/serverlib/api"
	_ "github.com/lib/pq"
)

func New(cfg api.Config) (*sql.DB, error) {
	fmt.Println("conn str ", cfg.DB.ConnStr)
	db, err := sql.Open("postgres", cfg.DB.ConnStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}