package data

import (
	"database/sql"
)

type Models struct {
	SQLScript interface {
		Register(script *SQLScript) error
		Get(name string) (*SQLScript, error)
		GetAll(projectName string) ([]*SQLScript, error)
	}
	SQLMigration interface {
		Add(mig *SQLMigration) error
		Update(mig *SQLMigration) error
		Remove(mig *SQLMigration) error
		GetAll(table string) ([]*SQLMigration, error)
		GetAllByDir(dir string, env string, table string) ([]*SQLMigration, error)
	}
	SQLMigrationGroup interface {
		Get(env string, table string) (*SQLMigrationGroup, error)
	}
	SQLMigrationsLatest interface {
		Update(mig *SQLMigrationsLatest) error
		AutoUpdateLatest(env string, table string) error
		GetByEnv(env string, table string) ([]*SQLMigrationsLatest, error)
	}
	SQLMigrationTables interface {
		Get(env string) (*SQLMigrationTables, error)
	}
	Project interface {
		Get(name string) (*Project, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		SQLScriptModel{DB: db},
		SQLMigrationModel{DB: db},
		SQLMigrationGroupModel{DB:db},
		SQLMigrationsLatestModel{DB: db},
		SQLMigrationTablesModel{DB: db},
		ProjectModel{DB: db},
	}
}
