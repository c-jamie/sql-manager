package migrate

import (
	"context"
	db "database/sql"
	"fmt"
	"strings"

	"github.com/c-jamie/sql-manager/clientlib/log"
	"github.com/rubenv/sql-migrate/sqlparse"
	
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
)

// Migration executes a migration 
type Migration struct {
	CTX context.Context
	DB  *db.DB
}

// New creates a new Migration struct
func New(driver string, connection string) (*Migration, error) {
	var mig Migration
	db, err := db.Open(driver, connection)
	if err != nil {
		return nil, err
	}
	mig.DB = db
	mig.DB.Ping()
	return &mig, nil
}

func (mig *Migration) parse(sql string) (*sqlparse.ParsedMigration, error) {

	migration, err := sqlparse.ParseMigration(strings.NewReader(sql))
	if err != nil {
		return nil, fmt.Errorf("migrate unable to build migration: %w", err)
	}
	return migration, nil
}

// Execute runs migrations in a provided direction
func (mig *Migration) Execute(sql string, dir string) error {

	if dir == "up" {
		migrations, err := mig.parse(sql)
		if err != nil {
			return err
		}
		for _, m := range migrations.UpStatements {
			log.Debug("executing up statement")
			err = mig.execute(m)
			if err != nil {
				return err
			}
		}
	}
	if dir == "down" {
		migrations, err := mig.parse(sql)
		if err != nil {
			return err
		}
		for _, m := range migrations.DownStatements {
			log.Debug("executing down statement")
			err = mig.execute(m)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CloseDB shuts the DB connection
func (mig *Migration) CloseDB() error {
	err := mig.DB.Close()
	return err
}

func (mig *Migration) execute(sql string) error {
	result, err := mig.DB.Exec(sql)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Debug(rows, " rows affected")
	return nil
}
