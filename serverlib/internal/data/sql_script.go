package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gosimple/slug"
)

type SQLScript struct {
	ProjectID       int    `json:"project_id"`
	Project         string `json:"project"`
	FileLocation    string `json:"file_location"`
	FileID          int    `json:"file_id"`
	SnippitID       int    `json:"snippit_id"`
	SnippitName     string `json:"snippit_name"`
	SnippitLocation string `json:"snippit_location"`
	Script          string `json:"script"`
}

type SQLScriptModel struct {
	DB *sql.DB
}

func (m SQLScriptModel) Register(script *SQLScript) error {
	query := `
		insert into project(name, created_at, updated_at)
		values 		($1, now(), now())
		returning 	id
	`
	args := []interface{}{script.Project}
	err := m.DB.QueryRow(query, args...).Scan(&script.ProjectID)

	if err != nil {
		return fmt.Errorf("unable to add project %w", err)
	}

	query = `
		insert into 	snippit(name, project_id, created_at, updated_at)
		values 			($1, $2, now(), now())
		returning 		id
	`

	args = []interface{}{slug.Make(script.FileLocation), script.ProjectID}
	err = m.DB.QueryRow(query, args...).Scan(&script.SnippitID)
	if err != nil {
		return fmt.Errorf("unable to add snippit %w", err)
	}

	query = `
		insert into 	git(location, snippit_id, created_at, updated_at)
		values 			($1, $2, now(), now())
		returning		id
	`
	args = []interface{}{script.FileLocation, script.SnippitID}
	err = m.DB.QueryRow(query, args...).Scan(&script.SnippitID)

	if err != nil {
		return fmt.Errorf("unable to add git %w", err)
	}
	return nil
}

func (m SQLScriptModel) Get(name string) (*SQLScript, error) {

	query := `
		select 		s.name
					, gt.location
		from 		snippit as s
		inner join	git as gt 
		on 			gt.snippit_id = s.id
		where		s.name = $1
	`

	var script SQLScript

	err := m.DB.QueryRow(query, name).Scan(&script.SnippitName, &script.FileLocation)

	if err != nil {
		return nil, err
	} else {
		return &script, nil
	}
}

func (m SQLScriptModel) GetAll(projectName string) ([]*SQLScript, error) {

	query := `
		select 		
					p.name as project_name
					, p.id as project_id
					, s.name as snippit_name
					, s.id as snippit_id
					, gt.location
					, gt.id as git_id
		from		project as p					
		inner join 	snippit as s
		on			p.id = s.project_id
		inner join	git as gt 
		on 			gt.snippit_id = s.id
		where		p.name = $1
	`

	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	rows, err := m.DB.QueryContext(ctx, query, projectName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var scripts []*SQLScript

	for rows.Next() {
		var script SQLScript

		err := rows.Scan(
			&script.Project,
			&script.ProjectID,
			&script.SnippitName,
			&script.SnippitID,
			&script.FileLocation,
			&script.FileID,
		)
		if err != nil {
			return nil, err
		}

		scripts = append(scripts, &script)
	}

	if err != nil {
		return nil, err
	} else {
		return scripts, nil
	}
}
