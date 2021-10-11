package api

import (
	"database/sql"
	"fmt"

	"github.com/c-jamie/sql-manager/serverlib/internal/data"
	"github.com/c-jamie/sql-manager/serverlib/internal/git"
	"github.com/c-jamie/sql-manager/serverlib/internal/migrations"
)

type Application struct {
	Config     *Config
	Models     data.Models
	Middleware Middleware
	Migrations migrations.Migrations
	GIT        git.Repo
}

type Config struct {
	Port        int
	Env         string
	Version     string
	GitURL      string
	GitUserName string
	GitToken    string
	Auth        string

	DB struct {
		ConnStr      string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdelTime  int
	}
}

func NewApplication(db *sql.DB, cfg *Config) (*Application, error) {
	git, err := git.New(cfg.GitURL, cfg.GitUserName, cfg.GitToken)
	if err != nil {
		return nil, fmt.Errorf("unable to start app %w", err)
	}
	app := Application{
		Config:     cfg,
		Models:     data.NewModels(db),
		GIT:        git,
		Middleware: NewMiddelware(cfg.Auth, db),
		Migrations: migrations.Migrations{DB: db},
	}

	return &app, nil
}
