package sql

import (
	"strings"
	"testing"

	"github.com/c-jamie/sql-manager/clientlib/app"
	sqlMig "github.com/c-jamie/sql-manager/clientlib/migration"
	"github.com/c-jamie/sql-manager/clientlib/mocks"
	"github.com/c-jamie/sql-manager/clientlib/utils"
	"github.com/c-jamie/sql-manager/clientlib/log"
)

func setup() {
	log.InitLog("info")

}

func setupApp() *app.App {
	clientApp := &app.App{
		ServerURL:  "/",
		CentralURL: "/",
		Home:       "/",
		AuthUrl:    "/",
		Token:      nil,
		Account:    nil,
		Script:     &mocks.Script{},
		Migration:  &mocks.Migration{},
	}
	return clientApp
}

func TestSQLLoadFile(t *testing.T) {
	setup()
	app := setupApp()
	app.Migration.Add("proj1/1_mig.sql", "dev", "a.b.c")
	testcases := []struct {
		sql string
		env string
		out string
	}{{
		sql: "../../resources/sql/test1.sql",
		env: "dev",
		out: "select 99",
	}}

	for _, test := range testcases {
		sqlScript, _ := utils.ReadFile(test.sql)
		mig, _ := app.Migration.GetAll(test.env)
		migEnv := make(map[string][]*sqlMig.SQLMigrationStrategy)
		migEnv[test.env] = mig
		sql := New(string(sqlScript), test.env, migEnv, app.Script.Get)
		parsed := sql.Compile()
		if !strings.Contains(parsed, test.out) {
			t.Error(parsed)
		}
	}
}

func TestSQLLoadServer(t *testing.T) {
	setup()
	app := setupApp()
	err := app.Migration.Add("proj1/1_mig.sql", "dev", "a.b.c")
	if err != nil {
		panic(err)
	}
	err = app.Script.Register("proj1/test3.sql")
	if err != nil {
		panic(err)
	}
	testcases := []struct {
		sql string
		env string
		out string
	}{{
		sql: "../../resources/sql/test5.sql",
		env: "dev",
		out: "A.B.test3",
	}}

	for _, test := range testcases {
		sqlScript, _ := utils.ReadFile(test.sql)
		mig, _ := app.Migration.GetAll(test.env)
		migEnv := make(map[string][]*sqlMig.SQLMigrationStrategy)
		migEnv[test.env] = mig
		sql := New(string(sqlScript), test.env, migEnv, app.Script.Get)
		parsed := sql.Compile()
		if !strings.Contains(parsed, test.out) {
			t.Error(parsed)
		}
	}
}
