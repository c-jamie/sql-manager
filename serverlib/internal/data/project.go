package data

import (
	"database/sql"
	"errors"
	"fmt"
)

type Project struct {
	Name       string       `json:"name"`
	ID         int          `json:"id"`
	SQLScripts []*SQLScript `json:"sql_script"`
}

type ProjectModel struct {
	DB *sql.DB
}

func (m ProjectModel) Get(name string) (*Project, error) {

	query := `
		select id, name from project where name = $1 and deleted_at is null
	`
	var project Project
	err := m.DB.QueryRow(query, name).Scan(&project.ID, &project.Name)

	sqlScriptModel := SQLScriptModel{DB: m.DB}

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("project does not exist %w", err)
		default:
			return nil, err
		}
	}

	sqlScripts, err := sqlScriptModel.GetAll(name)

	if err != nil {
		return nil, err
	}

	project.SQLScripts = sqlScripts
	return &project, nil
}
