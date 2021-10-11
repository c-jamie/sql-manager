package data

import (
	"database/sql"
	"fmt"
)

type SQLMigrationGroup struct {
	Env            string          `json:"env"`
	SourceTable    string          `json:"source_table"`
	MigrationsUp   []*SQLMigration `json:"migrations_up"`
	MigrationsDown []*SQLMigration `json:"migrations_down"`
}

type SQLMigrationGroupModel struct {
	DB *sql.DB
}

func (m SQLMigrationGroupModel) Get(env string, table string) (*SQLMigrationGroup, error) {

	migModel := SQLMigrationModel{DB: m.DB}
	var migGroup SQLMigrationGroup
	migUp, err := migModel.GetAllByDir("up", env, table)

	if err != nil {
		return nil, fmt.Errorf("unable to load up mig %w", err)
	}

	migDown, err := migModel.GetAllByDir("down", env, table)

	if err != nil {
		return nil, fmt.Errorf("unable to load up mig %w", err)
	}

	migGroup.MigrationsUp = migUp
	migGroup.MigrationsDown = migDown

	migGroup.Env = env
	migGroup.SourceTable = table

	return &migGroup, nil
}
