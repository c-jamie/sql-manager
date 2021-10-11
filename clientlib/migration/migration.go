package migation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/c-jamie/sql-manager/clientlib/log"
	"github.com/c-jamie/sql-manager/clientlib/migrate"
	"github.com/c-jamie/sql-manager/clientlib/request"
	"github.com/jedib0t/go-pretty/table"
	"github.com/tidwall/gjson"
)

const (
	AddMigration        = "migrations"
	DeleteMigration     = "migrations"
	GetMigration        = "migrations"
	ListMigrationTables = "migrations/table"
	SetMigration        = "migrations/latest"
	UpdateMigration     = "migrations"
)

// SQLMigrationStrategy defines the domain for a given table and env
// it includes our up and down migrations
type SQLMigrationStrategy struct {
	Table          string          `json:"source_table"`
	Env            string          `json:"env"`
	MigrationsUp   []*SQLMigration `json:"migrations_up"`
	MigrationsDown []*SQLMigration `json:"migrations_down"`
}

// SQLMigration defines the domain for a migration
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


// MigrationTable defines the domain for all migrations associated with a table
type MigrationTable struct {
	Table      string          `json:"Table"`
	Migrations []*SQLMigration `json:"Migrations"`
}


// Migration represents the interface used to control migrations
type Migration interface {
	// Add adds / registers a new migration
	Add(file string, env string, table string) error
	// Get returns the migration file
	Get(env string, table string) (*SQLMigrationStrategy, error)
	// ListTables lists the tables for a given env
	ListTables(env string) ([]*MigrationTable, error)
	// Delete removes a migration
	Delete(env string, table string, migrationID int) error
	// Update updates metadata associated with a migration
	Update(env string, table string, migrationID int, timeStamp time.Time, timeStampNull bool) error
	// Set defines the current latest migration
	Set(env string, table string, migrationID int) error
	// GetAll returns a all migrations for an env split by table
	GetAll(env string) ([]*SQLMigrationStrategy, error)
	// Run applies migrations
	Run(driver string, connection string, env string, table string) error
}

type migration struct {
	BaseURL string
}

// New creates a new migration container
func New(url string) Migration {
	mig := &migration{BaseURL: url}
	return mig
}

func (app *migration) makeRequest(url string, payload []byte, how string) ([]byte, error) {
	client := request.Client{BaseURL: app.BaseURL}
	var body []byte
	var resp *http.Response
	var err error
	if how == http.MethodPost || how == http.MethodPatch || how == http.MethodDelete {
		body, resp, err = client.PostReq(url, payload, true, how)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unable to make request %s", string(body))
		}
	} else {
		body, resp, err = client.GetReq(url, payload, true)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unable to make request %s", string(body))
		}
	}

	return body, nil
}

func (app *migration) List(env string) ([]*MigrationTable, error) {

	unmarshal := func(body []byte) ([]*MigrationTable, error) {
		var mig []*MigrationTable
		jsonErr := json.Unmarshal(body, &mig)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return mig, nil
	}
	payload := []byte(fmt.Sprintf(`{"env":"%s"}`, env))
	url := "/" + ListMigrationTables
	body, err := app.makeRequest(url, payload, http.MethodGet)
	if err != nil {
		return nil, err
	}
	out, err := unmarshal(body)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (app *migration) Add(file string, env string, table string) error {
	payload := []byte(fmt.Sprintf(`{"env":"%s", "table":"%s", "file": "%s"}`, env, table, file))
	url := "/" + AddMigration
	_, err := app.makeRequest(url, payload, http.MethodPost)
	if err != nil {
		return err
	}
	return nil
}

func (app *migration) Delete(env string, table string, migrationID int) error {
	payload := []byte(fmt.Sprintf(`{"env":"%s", "table":"%s", "sql_migration_id": %d}`, env, table, migrationID))
	url := "/" + DeleteMigration
	_, err := app.makeRequest(url, payload, http.MethodDelete)
	if err != nil {
		return err
	}
	return nil
}

func (app *migration) Get(env string, table string) (*SQLMigrationStrategy, error) {
	unmarshal := func(body []byte) (*SQLMigrationStrategy, error) {
		var mig SQLMigrationStrategy
		jsonErr := json.Unmarshal(body, &mig)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return &mig, nil
	}
	url := "/" + GetMigration + "?env=" + env + "&table=" + table
	body, err := app.makeRequest(url, nil, http.MethodGet)
	if err != nil {
		return nil, err
	}

	json := gjson.Get(string(body), "migrations").String()
	out, err := unmarshal([]byte(json))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (app *migration) ListTables(env string) ([]*MigrationTable, error) {
	unmarshal := func(body []byte) ([]*MigrationTable, error) {
		var mig []*MigrationTable
		jsonErr := json.Unmarshal(body, &mig)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return mig, nil
	}
	url := "/" + ListMigrationTables + "?env=" + env
	body, err := app.makeRequest(url, nil, http.MethodGet)
	json := gjson.Get(string(body), "migrations.tables").String()
	if err != nil {
		return nil, err
	}
	out, err := unmarshal([]byte(json))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (app *migration) Update(env string, table string, migrationID int, timeStamp time.Time, timeStampNull bool) error {
	payload := []byte(fmt.Sprintf(`{"env":"%s", "table":"%s", "sql_migration_id": %d, "migrated_at": "%s", "migrated_at_null": %t}`,
		env,
		table,
		migrationID,
		timeStamp,
		timeStampNull,
	))
	url := "/" + UpdateMigration
	_, err := app.makeRequest(url, payload, http.MethodPatch)
	if err != nil {
		return err
	}
	return nil
}

func (app *migration) Set(env string, table string, migrationID int) error {
	payload := []byte(fmt.Sprintf(`{"env":"%s", "table":"%s", "sql_migration_id": %d}`, env, table, migrationID))
	url := "/" + SetMigration
	_, err := app.makeRequest(url, payload, http.MethodPost)
	if err != nil {
		return err
	}
	return nil

}

func (app *migration) GetAll(env string) ([]*SQLMigrationStrategy, error) {
	var out []*SQLMigrationStrategy
	tables, err := app.ListTables(env)

	if err != nil {
		return nil, err
	}

	for _, t := range tables {
		mig, err := app.Get(env, t.Table)
		if err != nil {
			return nil, err
		}
		out = append(out, mig)
	}
	return out, nil
}

func (app *migration) Run(driver string, connection string, env string, table string) error {
	migrator, err := migrate.New(driver, connection)
	if err != nil {
		return err
	}
	migrations, err := app.Get(env, table)
	if err != nil {
		return err
	}
	log.Debug("loading migrations for ", driver, env, table)
	if migrations.MigrationsDown != nil {
		for _, sql := range migrations.MigrationsDown {
			log.Debug("migration down id", sql.ID, sql.MigratedAtNull)
			if !sql.MigratedAtNull {
				err = migrator.Execute(sql.Script, "down")
				if err != nil {
					return err
				}
				err = app.Update(env, table, sql.ID, time.Now(), true)
				if err != nil {
					return err
				}
			} else {
				log.Debug("migration ", sql.ID, " is already down")
			}
		}
	}
	if migrations.MigrationsUp != nil {
		for _, sql := range migrations.MigrationsUp {
			log.Debug("migration up id ", sql.ID, sql.MigratedAtNull)
			if sql.MigratedAtNull {
				err = migrator.Execute(sql.Script, "up")

				if err != nil {
					return err
				}
				err = app.Update(env, table, sql.ID, time.Now(), false)
				if err != nil {
					return err
				}
			} else {
				log.Debug("migration ", sql.ID, " already up")
			}
		}
	}
	return nil
}


// MigrationTablesToTable returns a text table of all migrations for an env
func MigrationTablesToTable(mit []*MigrationTable) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Table"})
	for i, k := range mit {
		t.AppendRow(table.Row{i, k.Table})
	}
	t.Render()
}


// MigrationsToTable returns a text table for a given migration strategy
func MigrationsToTable(mit SQLMigrationStrategy) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Dir", "ID", "File", "File Order", "Migrated At"})
	if mit.MigrationsUp != nil {
		for i, k := range mit.MigrationsUp {
			t.AppendRow(table.Row{i, "UP", k.ID, k.File, k.FileOrder, k.MigratedAt})
		}
	}
	if mit.MigrationsDown != nil {
		for i, k := range mit.MigrationsDown {
			t.AppendRow(table.Row{i, "DOWN", k.ID, k.File, k.FileOrder, k.MigratedAt})
		}
	}
	t.Render()
}
