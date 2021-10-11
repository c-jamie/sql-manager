package data

import (
	"fmt"
	"time"
	"context"
	"database/sql"
)

type SQLMigrationsLatest struct {
	Env             string `json:"env"`
	Table           string `json:"table"`
	SQLMigrationsID int    `json:"sql_migrations_id"`
}

type SQLMigrationsLatestModel struct {

	DB *sql.DB
}

func (m SQLMigrationsLatestModel) Update(mig *SQLMigrationsLatest) error {

	query := `
		update 		sql_migrations_latest
		set			sql_migrations_id 	= $1
		where		env 				= $2
		and			source_table 		= $3
	`
	args := []interface{}{mig.SQLMigrationsID, mig.Env, mig.Table}
	result, err := m.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("unable to update latest migration %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to update latest migration %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("unable to update latest migration")
	}

	return nil
}

func (m SQLMigrationsLatestModel) GetByEnv(env string, table string) ([]*SQLMigrationsLatest, error) {

	query := `
		select
					sml.source_table
					, sml.env
					, sml.sql_migrations_id
		from 		sql_migrations_latest as sml

		where 		sml.source_table = $1
		and 		sml.env 		 = $2
	`

	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	args := []interface{}{table, env}
	rows, err := m.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var migs []*SQLMigrationsLatest

	for rows.Next() {
		var mig SQLMigrationsLatest

		err := rows.Scan(
			&mig.Table,
			&mig.Env,
			&mig.SQLMigrationsID,
		)
		if err != nil {
			return nil, err
		}

		migs = append(migs, &mig)
	}

	if err != nil {
		return nil, err
	} else {
		return migs, nil
	}
}


func (m SQLMigrationsLatestModel) AutoUpdateLatest(env string, table string) (error) {
	query := `
		with sel as (
			select 		
						m.id
						, m.file_id
						, m.source_table
						, m.env
						, m.status
						, m.file
						, m.file_order
						, m.migrated_at
			from 		sql_migrations as m
			where		m.env	 		= $1
			and 		m.source_table 	= $2
			order by 	m.id desc
			limit 		1
		) 
		insert into sql_migrations_latest (sql_migrations_id, env, source_table)
		select id, env, source_table from sel	
	`
	args := []interface{}{env, table}
	result, err := m.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("unable to update latest migration %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to update latest migration %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("unable to update latest migration")
	}

	return nil
}