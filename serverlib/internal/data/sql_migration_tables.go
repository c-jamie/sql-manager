package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SQLMigrationTable struct {
	Table      string          `json:"table"`
	Migrations []*SQLMigration `json:"migrations"`
}

type SQLMigrationTables struct {
	Env    string               `json:"env"`
	Tables []*SQLMigrationTable `json:"tables"`
}

type SQLMigrationTablesModel struct {
	DB *sql.DB
}

func (m SQLMigrationTablesModel) Get(env string) (*SQLMigrationTables, error) {
	tables, err := m.tables(env)

	if err != nil {
		return nil, fmt.Errorf("unable to load migrations for env %w", err)
	}

	var migTables SQLMigrationTables
	migModel := SQLMigrationModel{DB: m.DB}
	for _, t := range tables {
		mig, err := migModel.GetAll(t)

		if err != nil {
			return nil, fmt.Errorf("unable to load migrations for env %w", err)
		}
		migTables.Tables = append(migTables.Tables, &SQLMigrationTable{Table: t, Migrations: mig})
	}

	migTables.Env = env

	return &migTables, nil

}

func (m SQLMigrationTablesModel) tables(env string) ([]string, error) {

	query := `
		select 	distinct source_table 
		from 	sql_migrations
		where 	env = $1
	`

	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	rows, err := m.DB.QueryContext(ctx, query, env)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tables []string

	for rows.Next() {
		var table string

		err := rows.Scan(
			&table,
		)
		if err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	if err != nil {
		return nil, err
	} else {
		return tables, nil
	}
}
