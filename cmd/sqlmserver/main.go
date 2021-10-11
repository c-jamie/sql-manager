package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/c-jamie/sql-manager/serverlib/api"
	"github.com/c-jamie/sql-manager/serverlib/db"
	"github.com/c-jamie/sql-manager/serverlib/log"
)

func main() {
	mode := os.Getenv("SQLM_SER_GIN_MODE")
	addr := os.Getenv("SQLM_SER_HTTP_PORT")
	debugLevel := "debug"
	if mode == "release" {
		debugLevel = "debug"
	}
	log.New(debugLevel)
	auth, ok := os.LookupEnv("SQLM_SER_AUTH")
	if !ok {
		log.Fatal("unable to load SQLM_SER_AUTH")
	}
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
	dbPort, ok := os.LookupEnv("SQLM_SER_DB_PORT")
	if !ok {
		log.Fatal("unable to load DB_PW")
	}
	ssl := "disable"

	log.Info("host is: ", dbHost)
	icn := os.Getenv("CLOUD_SQL_CONNECTION_NAME")
	dbConnStr := ""
	if icn != "" {
		dbHost = "/cloudsql/" + icn + "/"
		log.Info("dbhost: ", dbHost)
		dbConnStr = fmt.Sprintf("postgres://%s:%s@/%s?host=%s", dbUser, dbPW, dbName, dbHost)
	} else {
		dbConnStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPW, dbHost, dbPort, dbName, ssl)
	}
	log.Debug("db conn str ", dbConnStr)
	cfg := api.Config{}
	cfg.GitUserName = gitUserName
	cfg.GitToken = gitToken
	cfg.GitURL = gitUrl
	cfg.Auth = auth 
	cfg.Version = "v1"
	cfg.DB.ConnStr = dbConnStr
	cfg.DB.MaxIdelTime = 5
	cfg.DB.MaxOpenConns = 3

	db, err := db.New(cfg)
	if err != nil {
		log.Fatal("unable to start db: ", err)
	}

	app, err := api.NewApplication(db, &cfg)

	if err != nil {
		log.Fatal(err)
	}
	app.Migrations.DoMigrations("up")

	address := ":" + addr
	log.Info("listening on address: ", address)
	service := &http.Server{Addr: address, Handler: app.Routes(), ReadTimeout: 8 * time.Second, WriteTimeout: 8 * time.Second}
	log.Fatal(service.ListenAndServe())
}
