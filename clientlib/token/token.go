package token

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/c-jamie/sql-manager/clientlib/log"
	"github.com/c-jamie/sql-manager/clientlib/utils"
	"github.com/c-jamie/sql-manager/serverlib/api"
)

type Token interface {
	AddUser(ua *api.UserAccount) error
	AddGit(url string, tok string) (bool, error)
	AddToken(tok string) (error)
	Get() (*Client, error)
}

type token struct {
	Home    string
	BaseURL string
}

type Client struct {
	User  *api.UserAccount `json:"user"`
	Token string           `json:"token"`
	Cloud Cloud            `json:"cloud"`
	GIT   GIT              `json:"git"`
}

type Cloud struct {
	URL        string `json:"url"`
	LocationID string `json:"location_id"`
	ServiceID  string `json:"service_id"`
	ProjectID  string `json:"project_id"`
}

type GIT struct {
	GITUrl   string `json:"git_url"`
	GITToken string `json:"git_token"`
}

func New(url string, home string) Token {
	tok := &token{BaseURL: url, Home: home}
	_, err := utils.ReadFile(path.Join(home, "sqlmngr.json"))
	if err != nil {
		utils.ToFile(tok, path.Join(home, "sqlmngr.json"))
	}
	return tok
}

func (app *token) AddUser(ua *api.UserAccount) error {
	token := Client{}
	dat, err := utils.ReadFile(path.Join(app.Home, "sqlmngr.json"))
	if err != nil {
		return fmt.Errorf("unable to add provider %w", err)
	}
	err = json.Unmarshal(dat, &token)
	if err != nil {
		return fmt.Errorf("unable to add provider %w", err)
	}
	token.User = ua
	utils.ToFile(token, path.Join(app.Home, "sqlmngr.json"))
	return nil
}

func (app *token) AddGit(url string, tok string) (bool, error) {
	token := Client{}
	dat, err := utils.ReadFile(path.Join(app.Home, "sqlmngr.json"))
	if err != nil {
		return false, fmt.Errorf("unable to add git info %w", err)
	}
	err = json.Unmarshal(dat, &token)
	if err != nil {
		return false, fmt.Errorf("unable to add git info %w", err)
	}
	token.GIT = GIT{GITUrl: url, GITToken: tok}
	utils.ToFile(token, path.Join(app.Home, "sqlmngr.json"))
	return true, nil
}

func (app *token) AddToken(tok string) error {
	token := Client{}
	dat, err := utils.ReadFile(path.Join(app.Home, "sqlmngr.json"))
	if err != nil {
		return fmt.Errorf("unable to add git info %w", err)
	}
	err = json.Unmarshal(dat, &token)
	if err != nil {
		return fmt.Errorf("unable to add git info %w", err)
	}
	token.Token = tok
	utils.ToFile(token, path.Join(app.Home, "sqlmngr.json"))
	return nil
}

func (app *token) Get() (*Client, error) {
	token := Client{}
	dat, err := utils.ReadFile(path.Join(app.Home, "sqlmngr.json"))
	if err != nil {
		log.Debug(fmt.Errorf("unable to load token: %w", err))
		return nil, err
	}
	err = json.Unmarshal(dat, &token)
	if err != nil {
		log.Debug(fmt.Errorf("unable to load token: %w", err))
		return nil, err
	}
	return &token, nil
}
