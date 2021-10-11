package mocks

import (
	"time"

	mig "github.com/c-jamie/sql-manager/clientlib/migration"
)

type Migration struct {
}

func (app *Migration) Add(file string, env string, table string) error {
	return nil
}
func (app *Migration) Get(env string, table string) (*mig.SQLMigrationStrategy, error) {
	return nil, nil
}
func (app *Migration) ListTables(env string) ([]*mig.MigrationTable, error) {
	return nil, nil
}
func (app *Migration) Delete(env string, table string, migrationID int) error {
	return nil
}
func (app *Migration) Update(env string, table string, migrationID int, timeStamp time.Time, timeStampNull bool) error {
	return nil
}
func (app *Migration) Set(env string, table string, migrationID int) error {
	return nil
}
func (app *Migration) GetAll(env string) ([]*mig.SQLMigrationStrategy, error) {

	strategy := []*mig.SQLMigrationStrategy{{Table: "a.b.c", Env: "int", MigrationsUp: []*mig.SQLMigration{{ID: 1}}, MigrationsDown: []*mig.SQLMigration{{ID: 1}}}}

	return strategy, nil
}
func (app *Migration) Run(driver string, connection string, env string, table string) error {
	return nil
}
