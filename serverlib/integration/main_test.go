package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/c-jamie/sql-manager/serverlib/api"
	"github.com/c-jamie/sql-manager/serverlib/db"
	"github.com/c-jamie/sql-manager/serverlib/internal/data"
	"github.com/c-jamie/sql-manager/serverlib/internal/git"
	"github.com/c-jamie/sql-manager/serverlib/log"
	"github.com/c-jamie/sql-manager/serverlib/internal/migrations"
	"github.com/c-jamie/sql-manager/serverlib/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type ClientToken struct {
	Email string    `json:"email"`
	Token AuthToken `json:"token"`
}

func DoRequest(app *api.Application, json []byte, url string, token string, method string) (*bytes.Buffer, int) {
	testRouter := app.Routes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(json))
	req.Header.Set("Authorization", "Bearer "+token)
	testRouter.ServeHTTP(w, req)
	return w.Body, w.Code
}

func setup() *api.Application {
	log.New("debug")

	gitUserName, ok := os.LookupEnv("SQLM_SER_GIT_USERNAME")
	if !ok {
		log.Fatal("unable to load SQLM_SER_GIT_USERNAME")
	}
	gitToken, ok := os.LookupEnv("SQLM_SER_GIT_TOKEN")
	if !ok {
		log.Fatal("unable to load SQLM_SER_GIT_TOKEN")
	}
	gitUrl, ok := os.LookupEnv("SQLM_SER_GIT_URL")
	if !ok {
		log.Fatal("unable to load SQLM_SER_GIT_URL")
	}
	dbHost, ok := os.LookupEnv("SQLM_SER_DB_HOST")
	if !ok {
		log.Fatal("unable to load DB_HOST")
	}
	dbPort, ok := os.LookupEnv("SQLM_SER_DB_PORT")
	if !ok {
		log.Fatal("unable to load DB_HOST")
	}
	dbName, ok := os.LookupEnv("SQLM_SER_DB_NAME")
	if !ok {
		log.Fatal("unable to load DB_NAME")
	}
	dbUser, ok := os.LookupEnv("SQLM_SER_DB_USER")
	if !ok {
		log.Fatal("unable to load DB_USER")
	}
	dbPW, ok := os.LookupEnv("SQLM_SER_DB_PW")
	if !ok {
		log.Fatal("unable to load DB_PW")
	}
	if !ok {
		log.Fatal("unable to load DB_PW")
	}
	ssl := "disable"
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPW, dbHost, dbPort, dbName, ssl)
	log.Info("db conn str is ", dbConnStr)

	cfg := api.Config{}
	cfg.Version = "v1"
	cfg.DB.ConnStr = dbConnStr
	cfg.DB.MaxIdelTime = 5
	cfg.DB.MaxOpenConns = 3

	db, err := db.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	git, err := git.New(gitUrl, gitUserName, gitToken)

	if err != nil {
		log.Fatal(err)
	}

	app := api.Application{
		Config:     &cfg,
		Models:     data.NewModels(db),
		Middleware: &mocks.MockMiddleware{},
		Migrations: migrations.Migrations{DB: db},
		GIT:        git,
	}
	app.Migrations.DoMigrations("up")

	return &app
}

func GetAuthToken() ClientToken {
	token := ClientToken{}
	dat, _ := ioutil.ReadFile(path.Join("/home/cillairne/code/polar", "sqlmngr.json"))
	err := json.Unmarshal(dat, &token)
	if err != nil {
		log.Fatal("unable to load token: %s", err)
	}
	return token
}

func TestPingRoute(t *testing.T) {
	app := setup()
	bearer := GetAuthToken()
	testRouter := app.Routes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	req.Header.Set("Authorization", "Bearer "+bearer.Token.IDToken)
	testRouter.ServeHTTP(w, req)
	t.Log(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestScriptRegister(t *testing.T) {
	testcases := []struct {
		in     []byte
		code   int
		expect string
	}{
		{
			in:     []byte(fmt.Sprintf(`{"file":"%s", "project":"%s"}`, "proj1/test1.sql", "pro-one")),
			code:   http.StatusCreated,
			expect: "files",
		}, {
			in:     []byte(fmt.Sprintf(`{"file":"%s", "project":"%s"}`, "proj1/test1.sql", "")),
			code:   http.StatusUnprocessableEntity,
			expect: "must not be empty",
		},
	}
	app := setup()
	bearer := GetAuthToken()

	for _, tcase := range testcases {
		out, code := DoRequest(app, tcase.in, "/v1/files", bearer.Token.AccessToken, http.MethodPost)
		t.Log(out.String())
		assert.Equal(t, tcase.code, code)
		assert.Equal(t, strings.Contains(out.String(), tcase.expect), true)
	}
	app.Migrations.DoMigrations("down")
}
func TestScriptGet(t *testing.T) {
	testcases := []struct {
		in     string
		model  data.SQLScript
		code   int
		expect string
	}{
		{
			in:     "/v1/files?name=proj1-test1-sql",
			model:  data.SQLScript{FileLocation: "proj1/test1.sql", Project: "test1"},
			code:   http.StatusOK,
			expect: "files",
		}, {
			in:     "/v1/files",
			model:  data.SQLScript{FileLocation: "proj1/test2.sql", Project: "test1"},
			code:   http.StatusUnprocessableEntity,
			expect: "must not be empty",
		},
	}
	app := setup()
	bearer := GetAuthToken()

	for _, tcase := range testcases {
		app.Models.SQLScript.Register(&tcase.model)
		out, code := DoRequest(app, []byte(""), tcase.in, bearer.Token.AccessToken, http.MethodGet)
		assert.Equal(t, tcase.code, code)
		assert.Equal(t, strings.Contains(out.String(), tcase.expect), true)
	}
	app.Migrations.DoMigrations("down")
}

func TestProjectGet(t *testing.T) {
	testcases := []struct {
		in     string
		model  []data.SQLScript
		code   int
		expect string
	}{
		{
			in: "/v1/projects?name=test1",
			model: []data.SQLScript{
				{FileLocation: "proj1/test1.sql", Project: "test1"},
				{FileLocation: "proj1/test2.sql", Project: "test1"},
			},
			code:   http.StatusOK,
			expect: "projects",
		}, {
			in: "/v1/projects?name=",
			model: []data.SQLScript{
				{FileLocation: "proj1/test2.sql", Project: "test1"},
			},
			code:   http.StatusUnprocessableEntity,
			expect: "must not be empty",
		},
	}
	app := setup()
	bearer := GetAuthToken()

	for _, tcase := range testcases {
		for _, m := range tcase.model {
			app.Models.SQLScript.Register(&m)
		}
		out, code := DoRequest(app, []byte(""), tcase.in, bearer.Token.AccessToken, "GET")
		t.Log(out.String())
		assert.Equal(t, tcase.code, code)
		assert.Equal(t, strings.Contains(out.String(), tcase.expect), true)
	}
	app.Migrations.DoMigrations("down")
}

func TestFilesListAll(t *testing.T) {
	testcases := []struct {
		in     string
		model  []data.SQLScript
		code   int
		expect string
	}{
		{
			in: "/v1/files/list?project=test1",
			model: []data.SQLScript{
				{FileLocation: "proj1/test1.sql", Project: "test1"},
				{FileLocation: "proj1/test2.sql", Project: "test1"},
			},
			code:   http.StatusOK,
			expect: "files",
		},
	}
	app := setup()
	bearer := GetAuthToken()

	for _, tcase := range testcases {
		for _, m := range tcase.model {
			app.Models.SQLScript.Register(&m)
		}
		out, code := DoRequest(app, []byte(""), tcase.in, bearer.Token.AccessToken, "GET")
		t.Log(out.String())
		assert.Equal(t, tcase.code, code)
		assert.Equal(t, strings.Contains(out.String(), tcase.expect), true)
	}
	app.Migrations.DoMigrations("down")
}

func TestAddMigration(t *testing.T) {
	testcases := []struct {
		in     []byte
		code   int
		expect string
	}{
		{
			in:     []byte(`{"file":"dir1/dir2/1_init.sql", "env":"dev", "table":"db.sch.table1"}`),
			code:   http.StatusCreated,
			expect: "dir1-dir2-1_init-sql",
		}, {
			in:     []byte(`{"file":"dir1/dir2/2_init.sql", "env":"dev", "table":"db.sch.table1"}`),
			code:   http.StatusCreated,
			expect: "dir1-dir2-2_init-sql",
		}, {
			in:     []byte(`{"file":"dir1/dir2/2_init.sql", "env":"dev", "table":"db.sch.table2"}`),
			code:   http.StatusCreated,
			expect: "dir1-dir2-2_init-sql",
		},
	}
	app := setup()
	bearer := GetAuthToken()
	for _, tcase := range testcases {
		out, code := DoRequest(app, tcase.in, "/v1/migrations", bearer.Token.AccessToken, http.MethodPost)
		assert.Equal(t, tcase.code, code)
		assert.Equal(t, tcase.expect, gjson.Get(out.String(), "migrations.file_id").Str, tcase.expect)
		t.Log(out.String())
	}
	app.Migrations.DoMigrations("down")
}

func TestGetMigration(t *testing.T) {
	testcases := []struct {
		file   string
		env    string
		table  string
		code   int
		expect string
		url    string
	}{
		{
			file:  "dir1/dir2/1_init.sql",
			env:   "deb",
			table: "db.sch.table1",
			code:  http.StatusOK,
			url:   "/v1/migrations?env=deb&table=db.sch.table1",
		},
	}
	app := setup()
	bearer := GetAuthToken()
	for _, tcase := range testcases {
		mig := data.SQLMigration{File: tcase.file, Env: tcase.env, SourceTable: tcase.table}
		err := app.Models.SQLMigration.Add(&mig)
		assert.Equal(t, err, nil)
		out, code := DoRequest(app, []byte(""), tcase.url, bearer.Token.AccessToken, http.MethodGet)
		assert.Equal(t, tcase.code, code)
		assert.Equal(t, tcase.table, gjson.Get(out.String(), "migrations.source_table").Str)
		t.Log(out.String())
	}
	app.Migrations.DoMigrations("down")
}

func TestDeleteMigration(t *testing.T) {
	testcases := []struct {
		in     []byte
		code   int
		url    string
		id     int64
		models []data.SQLMigration
	}{
		{
			in:   []byte(`{"sql_migration_id":1, "env":"dev", "table":"db.sch.tb1"}`),
			code: http.StatusOK,
			url:  "/v1/migrations",
			id:   1,
			models: []data.SQLMigration{
				{File: "dir1/dir2/1_init.sql", Env: "dev", SourceTable: "db.sch.tb1"},
				{File: "dir1/dir2/2_init.sql", Env: "dev", SourceTable: "db.sch.tb1"},
			},
		},
	}
	app := setup()
	bearer := GetAuthToken()
	for _, tcase := range testcases {
		for _, m := range tcase.models {
			err := app.Models.SQLMigration.Add(&m)
			assert.Equal(t, err, nil)
		}
		out, code := DoRequest(app, tcase.in, tcase.url, bearer.Token.AccessToken, http.MethodDelete)
		assert.Equal(t, tcase.code, code)
		assert.Equal(t, tcase.id, gjson.Get(out.String(), "migrations.id").Int())
		latest, err :=  app.Models.SQLMigrationsLatest.GetByEnv("dev", "db.sch.tb1")
		assert.Equal(t, err, nil)
		t.Log(len(latest))
		t.Log(out.String())
	}
	app.Migrations.DoMigrations("down")
}

func TestSetLatestMigration(t *testing.T) {
	testcases := []struct {
		in            []byte
		code          int
		url           string
		migrations_id int64
		models        []data.SQLMigration
	}{
		{
			in:            []byte(`{"sql_migration_id":1, "env":"dev", "table":"db.sch.tb1"}`),
			code:          http.StatusOK,
			url:           "/v1/migrations/latest",
			migrations_id: 1,
			models: []data.SQLMigration{
				{File: "dir1/dir2/1_init.sql", Env: "dev", SourceTable: "db.sch.tb1"},
				{File: "dir1/dir2/2_init.sql", Env: "dev", SourceTable: "db.sch.tb1"}},
		},
	}
	app := setup()
	bearer := GetAuthToken()
	for _, tcase := range testcases {
		for _, m := range tcase.models {
			err := app.Models.SQLMigration.Add(&m)
			assert.Equal(t, err, nil)
		}
		out, code := DoRequest(app, tcase.in, tcase.url, bearer.Token.AccessToken, http.MethodPost)
		assert.Equal(t, tcase.code, code)
		assert.Equal(t, tcase.migrations_id, gjson.Get(out.String(), "migrations_latest.sql_migrations_id").Int())
		t.Log(out.String())
	}
	app.Migrations.DoMigrations("down")
}

func TestListMigrationTables(t *testing.T) {
	testcases := []struct {
		in            []byte
		code          int
		url           string
		migrations_id int64
		models        []data.SQLMigration
	}{
		{
			in:            []byte(""),
			code:          http.StatusOK,
			url:           "/v1/migrations/table?env=dev",
			migrations_id: 1,
			models: []data.SQLMigration{
				{File: "dir1/dir2/1_init.sql", Env: "dev", SourceTable: "db.sch.tb1"},
				{File: "dir1/dir2/2_init.sql", Env: "dev", SourceTable: "db.sch.tb2"},
				{File: "dir1/dir2/3_init.sql", Env: "dev", SourceTable: "db.sch.tb1"},
				{File: "dir1/dir2/4_init.sql", Env: "dev", SourceTable: "db.sch.tb2"},
			},
		},
	}
	app := setup()
	bearer := GetAuthToken()
	for _, tcase := range testcases {
		for _, m := range tcase.models {
			err := app.Models.SQLMigration.Add(&m)
			assert.Equal(t, err, nil)
		}
		out, code := DoRequest(app, tcase.in, tcase.url, bearer.Token.AccessToken, http.MethodGet)
		assert.Equal(t, tcase.code, code)
		t.Log(out.String())
	}
	app.Migrations.DoMigrations("down")
}
func TestUpdateMigration(t *testing.T) {
	testcases := []struct {
		in            []byte
		code          int
		url           string
		migrations_id int64
		models        []data.SQLMigration
	}{
		{
			in:            []byte(`{"sql_migration_id":3, "env":"dev", "table":"db.sch.tb1", "migrated_at_null": true}`),
			code:          http.StatusOK,
			url:           "/v1/migrations",
			migrations_id: 1,
			models: []data.SQLMigration{
				{File: "dir1/dir2/1_init.sql", Env: "dev", SourceTable: "db.sch.tb1"},
				{File: "dir1/dir2/2_init.sql", Env: "dev", SourceTable: "db.sch.tb2"},
				{File: "dir1/dir2/3_init.sql", Env: "dev", SourceTable: "db.sch.tb1"},
				{File: "dir1/dir2/4_init.sql", Env: "dev", SourceTable: "db.sch.tb2"},
			},
		},
	}
	app := setup()
	bearer := GetAuthToken()
	for _, tcase := range testcases {
		for _, m := range tcase.models {
			err := app.Models.SQLMigration.Add(&m)
			assert.Equal(t, err, nil)
		}
		out, code := DoRequest(app, tcase.in, tcase.url, bearer.Token.AccessToken, http.MethodPatch)
		assert.Equal(t, tcase.code, code)
		t.Log(out.String())
	}
	app.Migrations.DoMigrations("down")
}
