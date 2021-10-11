package app

import (
	"crypto/rand"
	_ "embed"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"github.com/c-jamie/sql-manager/clientlib/account"
	"github.com/c-jamie/sql-manager/clientlib/log"
	migation "github.com/c-jamie/sql-manager/clientlib/migration"
	"github.com/c-jamie/sql-manager/clientlib/request"
	"github.com/c-jamie/sql-manager/clientlib/script"
	"github.com/c-jamie/sql-manager/clientlib/token"

)

var CLIENT_MODE string

// App represents our application
type App struct {
	ServerURL  string
	CentralURL string
	AppVersion string
	AuthUrl    string
	Home       string
	Token      token.Token
	Account    account.Account
	Migration  migation.Migration
	Script     script.Script
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// New returns a new App
func New(level string) (*App, error) {
	log := log.InitLog(level)
	server := getEnv("SQL_MNGR_SERVER", "")
	authURL := getEnv("SQL_MNGR_AUTH", "")
	home := getEnv("SQL_MNGR_HOME", "")
	urlVersion := "v1"
	CLIENT_MODE = getEnv("SQL_MNGR_MODE", "")
	var err error
	if authURL == "" {
		return nil, fmt.Errorf("SQL_MNGR_AUTH is not set")
	}
	if home == "" {
		home, err = os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("unable to find the home dir")
		}
	}

	tok := token.New((server + "/" + urlVersion), home)
	scr := script.New((server + "/" + urlVersion))
	mig := migation.New((server + "/" + urlVersion))
	acc := account.New(authURL + "/" + urlVersion, home)

	app := &App{
		ServerURL:  (server + "/" + urlVersion),
		Home:       home,
		AuthUrl:    authURL + "/" + urlVersion,
		Token:      tok,
		Account:    acc,
		Script:     scr,
		Migration:  mig,
	}
	token, err := tok.Get()
	if err != nil {
		return app, nil
	} else {
		if app.ServerURL == "" {
			log.Debug("loading server URL from token:", token.Cloud.URL)
			app.ServerURL = token.Cloud.URL + "/" + urlVersion
		}
		return app, nil
	}
}

func (app *App) makeRequest(url string, payload []byte, how string) ([]byte, error) {
	client := request.Client{BaseURL: app.ServerURL}
	var body []byte
	var resp *http.Response
	if how == http.MethodPost {
		body, resp, _ = client.PostReq(url, payload, true, how)
		if resp.StatusCode != http.StatusCreated {
			return nil, fmt.Errorf("unable to make request")
		}
	} else {
		body, resp, _ = client.GetReq(url, payload, true)
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unable to make request")
		}

	}

	return body, nil
}

func (app *App) GetVersion() (string, error) {
	return "", nil
}

func (app *App) Server() (string, error) {
	url := "/healthcheck"
	body, err := app.makeRequest(url, nil, http.MethodGet)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	return string(body), nil
}

func (app *App) GenerateRandomString(n int) (string, error) {
	const letters = "C123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
