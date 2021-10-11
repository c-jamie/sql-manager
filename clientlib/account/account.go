package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/c-jamie/sql-manager/clientlib/log"
	"github.com/c-jamie/sql-manager/clientlib/request"
	cutils "github.com/c-jamie/sql-manager/clientlib/utils"
	"github.com/c-jamie/sql-manager/serverlib/api"
	"github.com/jedib0t/go-pretty/table"
	"github.com/tidwall/gjson"
)

// Account is the interface used to maintain the state of a given team.
type Account interface {
	// Register is used to create a new user on the platform.
	Register(email string, team string, password string) (*api.UserAccount, error)
	// Login is used to login to platform and get a token.
	Login(email string, password string) (string, error)
	// RemoveUser is used to delete a user from the platform.
	RemoveUser(user_id int) error
	// User is used to return the current user.
	User() (*api.UserAccount, error)
}

type account struct {
	BaseURL string
	Home    string
}

const (
	Users = "users"
	token = "tokens/authentication"
)

func New(url string, home string) Account {
	acc := &account{BaseURL: url, Home: home}
	return acc
}

func (app *account) makeRequest(url string, payload []byte, how string) ([]byte, error) {
	client := request.Client{BaseURL: app.BaseURL}
	var body []byte
	var resp *http.Response
	if how == http.MethodPost {
		body, resp, _ = client.PostReq(url, payload, true, how)
		if resp.StatusCode != http.StatusCreated {
			return nil, fmt.Errorf("%s", string(body))
		}
	} else {
		body, resp, _ = client.GetReq(url, payload, true)
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s", string(body))
		}

	}

	return body, nil
}

func (app *account) teamID() (int64, error) {
	dat, err := cutils.ReadFile(path.Join(app.Home, "sqlmngr.json"))

	if err != nil {
		log.Debug(fmt.Errorf("unable to load token: %w", err))
		return -1, err
	}
	token := gjson.Get(string(dat), "team_id").Int()
	return token, nil
}

func (app *account) Register(email string, team string, password string) (*api.UserAccount, error) {
	url := "/" + Users
	unmarshal := func(body []byte) (*api.UserAccount, error) {
		var ua *api.UserAccount
		jsonErr := json.Unmarshal(body, &ua)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return ua, nil
	}
	var jsonStr = []byte(fmt.Sprintf(`{"team": "%s", "email":"%s", "password": "%s"}`, team, email, password))
	resp, err := app.makeRequest(url, jsonStr, http.MethodPost)
	if err != nil {
		return nil, err
	}
	userJson := gjson.Get(string(resp), "users")
	if err != nil {
		return nil, err
	}
	user, err := unmarshal([]byte(userJson.String()))
	if err != nil {
		return nil, err
	}
	return user, err
}

func (app *account) Login(email string, password string) (string, error) {
	url := "/" + token
	var jsonStr = []byte(fmt.Sprintf(`{"email":"%s", "password": "%s"}`, email, password))
	resp, err := app.makeRequest(url, jsonStr, http.MethodPost)
	if err != nil {
		return "", err
	}
	token := gjson.Get(string(resp), "authentication_token.plain_text").String()
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

func (app *account) RemoveUser(user_id int) error {
	url := "/" + Users
	teamID, err := app.teamID()
	if err != nil {
		return err
	}
	var jsonStr = []byte(fmt.Sprintf(`{"user_id": %d, "team_id": %d}`, user_id, teamID))
	body, err := app.makeRequest(url, jsonStr, http.MethodPost)
	if err != nil {
		return fmt.Errorf("unable to register team %s", body)
	} else {
		return nil
	}
}

func (app *account) User() (*api.UserAccount, error) {
	unmarshal := func(body []byte) (*api.UserAccount, error) {
		var ua *api.UserAccount
		jsonErr := json.Unmarshal(body, &ua)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return ua, nil
	}
	url := "/" + Users
	body, err := app.makeRequest(url,nil, http.MethodGet)
	if err != nil {
		return nil, fmt.Errorf("request unavailable %s", body)
	}
	userJson := gjson.Get(string(body), "users")
	users, err := unmarshal([]byte(userJson.String()))
	if err != nil {
		return nil, fmt.Errorf("unable to load user account %w", err)
	} else {
		return users, nil
	}
}

func UsersToTable(users []*api.UserAccount) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"User ID", "Email"})
	for _, k := range users {
		t.AppendRow(table.Row{
			k.ID, k.Email,
		})
	}
	t.Render()
}
