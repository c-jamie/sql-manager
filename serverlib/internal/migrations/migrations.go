package migrations

import (
	"database/sql"
	"embed"
	"net/http"

	"github.com/c-jamie/sql-manager/serverlib/log"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed *.sql
var Files embed.FS

type Migrations struct {
	DB *sql.DB
}

func (app *Migrations) DoMigrations(dir string) {
	migrationSource := &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(Files),
	}
	if dir == "up" {
		n, err := migrate.Exec(app.DB, "postgres", migrationSource, migrate.Up)
		if err != nil {
			log.Fatal("migrations failed:", err)
		}
		log.Info("migrations up completed ", n)
	} else {
		n, err := migrate.Exec(app.DB, "postgres", migrationSource, migrate.Down)
		if err != nil {
			log.Fatal("migrations failed:", err)
		}
		log.Info("migrations down completed ", n)
	}
}
