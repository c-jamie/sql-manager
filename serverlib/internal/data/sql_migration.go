package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gosimple/slug"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type SQLMigration struct {
	Env            string     `json:"env"`
	FileID         string     `json:"file_id"`
	File           string     `json:"file"`
	FileOrder      int        `json:"file_order"`
	SourceTable    string     `json:"source_table"`
	MigratedAt     *time.Time `json:"migrated_at"`
	MigratedAtNull bool       `json:"migrated_at_null"`
	ID             int        `json:"id"`
	Script         string     `json:"script"`
	Status         *string    `json:"status"`
}

type SQLMigrationModel struct {
	DB *sql.DB
}

func (m SQLMigrationModel) Add(mig *SQLMigration) error {

	fileid := slug.Make(mig.File)
	fileid = strings.Join(strings.Split(fileid, "."), "-")
	mig.FileID = fileid
	file := filepath.Base(mig.File)
	order := strings.Split(file, "_")
	if len(order) == 0 {
		return fmt.Errorf("unable to add migration")
	}
	orderInt, err := strconv.Atoi(order[0])
	if err != nil {
		return fmt.Errorf("unable to add migration %w", err)
	}

	query := `
		insert into sql_migrations(env, file_id, file, file_order, source_table, migrated_at)
		values 		($1, $2, $3, $4, $5, null)
		returning 	id
	`
	args := []interface{}{mig.Env, fileid, mig.File, orderInt, mig.SourceTable}
	err = m.DB.QueryRow(query, args...).Scan(&mig.ID)

	if err != nil {
		return fmt.Errorf("unable to add project 1 %w", err)
	}

	query = `
		select id from sql_migrations_latest where source_table = $1 and env = $2
	`
	args = []interface{}{mig.SourceTable, mig.Env}
	result, err := m.DB.Exec(query, args...)

	if err != nil {
		return fmt.Errorf("unable to add project 2 %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to add project 3 %w", err)
	}

	if affected == 0 {
		query = `
			insert into sql_migrations_latest(sql_migrations_id, source_table, env)
			values 		($1, $2, $3) 
			returning 	id
		`
		args = []interface{}{mig.ID, mig.SourceTable, mig.Env}
		err = m.DB.QueryRow(query, args...).Scan(&mig.ID)
		if err != nil {
			return fmt.Errorf("unable to add project 4 %w", err)
		}
	} else {
		query = `
			update 		sql_migrations_latest
			set			sql_migrations_id 	= $1
			where		env 				= $2
			and			source_table 		= $3
		`
		args = []interface{}{mig.ID, mig.Env, mig.SourceTable}
		result, err = m.DB.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("unable to add project 5 %w", err)
		}
		rows, err := result.RowsAffected()

		if err != nil {
			return fmt.Errorf("unable to add project 5 %w", err)
		}

		if rows == 0 {
			return fmt.Errorf("unable to add project %w", sql.ErrNoRows)
		}

	}
	return nil

}
func (m SQLMigrationModel) Update(mig *SQLMigration) error {
	query := ""

	if mig.MigratedAtNull {
		query = `
			update 		sql_migrations
			set			migrated_at 	= null
			where		id 				= $1
			and			env 			= $2
			returning 	migrated_at
		`
	} else {
		query = `
			update 		sql_migrations
			set			migrated_at 	= now()
			where		id 				= $1
			and			env 			= $2
			returning 	migrated_at
		`

	}
	args := []interface{}{mig.ID, mig.Env}
	err := m.DB.QueryRow(query, args...).Scan(&mig.MigratedAt)

	if err != nil {
		return fmt.Errorf("unable to update project %w", err)
	}

	return nil
}

func (m SQLMigrationModel) Remove(mig *SQLMigration) error {
	query := `
		delete from sql_migrations 
		where		id = $1
	`
	args := []interface{}{mig.ID}
	result, err := m.DB.Exec(query, args...)

	if err != nil {
		return fmt.Errorf("unable to remove project %w", err)
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("unable to remove project %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no rows affected")
	}

	migLatest := SQLMigrationsLatestModel{DB: m.DB}
	latestMigrations, err := migLatest.GetByEnv(mig.Env, mig.SourceTable)

	if err != nil {
		return fmt.Errorf("unable to remove project %w", err)
	}

	if len(latestMigrations) == 0 {
		
		err := migLatest.AutoUpdateLatest(mig.Env, mig.SourceTable)
		if err != nil {
			return fmt.Errorf("unable to remove project %w", err)
		}

	}

	return nil
}

func (m SQLMigrationModel) Get(id int) (*SQLMigration, error) {

	query := `
		select 		m.id
					, m.file_id
					, m.source_table
					, m.env
					, m.status
					, m.file
					, m.file_order
					, m.migrated_at
					, case when m.migrated_at is null
					then true
					else false
					end as migrated_at_null
		from 		sql_migrations as m
		where		m.id = $1
	`

	var mig SQLMigration

	err := m.DB.QueryRow(query, id).Scan(
		&mig.ID, 
		&mig.FileID, 
		&mig.SourceTable, 
		&mig.Env, 
		&mig.Status, 
		&mig.File, 
		&mig.FileOrder, 
		&mig.MigratedAt,
		&mig.MigratedAtNull,
	)

	if err != nil {
		return nil, err
	} else {
		return &mig, nil
	}
}

func (m SQLMigrationModel) GetAll(table string) ([]*SQLMigration, error) {
	query := `
		select 		m.id
					, m.file_id
					, m.source_table
					, m.env
					, m.status
					, m.file
					, m.file_order
					, m.migrated_at
		from 		sql_migrations as m
		where 		m.source_table = $1
	`

	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	rows, err := m.DB.QueryContext(ctx, query, table)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var migs []*SQLMigration

	for rows.Next() {
		var mig SQLMigration

		err := rows.Scan(
			&mig.ID,
			&mig.FileID,
			&mig.SourceTable,
			&mig.Env,
			&mig.Status,
			&mig.File,
			&mig.FileOrder,
			&mig.MigratedAt,
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

func (m SQLMigrationModel) GetAllByDir(dir string, env string, table string) ([]*SQLMigration, error) {

	fileEQ := "<="
	if dir == "down" {
		fileEQ = ">"
	}
	query := fmt.Sprintf(`
		with
			sql_latest_order as (
			select
						sml.source_table
						, sml.env
						, sm.file_order
			from 		sql_migrations as sm
			inner join 	sql_migrations_latest as sml

			on 			sm.source_table = sml.source_table
			and 		sm.id 			= sml.sql_migrations_id
			where 		sm.source_table = $1
			and 		sm.env 			= $2
		)
		select 		m.id
					, m.file_id
					, m.source_table
					, m.env
					, m.status
					, m.file
					, m.file_order
					, m.migrated_at
					, case when m.migrated_at is null
					then true
					else false
					end as migrated_at_null
		from 		sql_migrations as m

		inner join 	sql_latest_order as sml 
		on 			m.source_table	= sml.source_table
		and			m.env 			= sml.env
		Where 		m.file_order %s sml.file_order
		Order by 	m.file_order asc
	`, fileEQ)

	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	rows, err := m.DB.QueryContext(ctx, query, table, env)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var migs []*SQLMigration

	for rows.Next() {
		var mig SQLMigration

		err := rows.Scan(
			&mig.ID,
			&mig.FileID,
			&mig.SourceTable,
			&mig.Env,
			&mig.Status,
			&mig.File,
			&mig.FileOrder,
			&mig.MigratedAt,
			&mig.MigratedAtNull,
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
