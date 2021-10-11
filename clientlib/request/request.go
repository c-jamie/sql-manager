package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/c-jamie/sql-manager/clientlib/log"
	"github.com/c-jamie/sql-manager/clientlib/utils"
	"github.com/tidwall/gjson"
)


// Client represents a http Client
type Client struct {
	BaseURL    string
	Home       string
	httpClient http.Client
}

func (app *Client) token() (string, error) {
	dat, err := utils.ReadFile(path.Join(app.Home, "sqlmngr.json"))

	if err != nil {
		log.Debug(fmt.Errorf("unable to load token: %w", err))
		return "", err
	}
	token := gjson.Get(string(dat), "token").String()
	return token, nil
}


// GetReq performs a http GET request on a resource
func (client *Client) GetReq(url string, jsonStr []byte, token bool) ([]byte, *http.Response, error) {
	url = client.BaseURL + url
	log.Debug("calling: ", url)
	log.Debug("payload: ", string(jsonStr))
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, nil, fmt.Errorf("server req failed %w", err)
	}
	req.Header.Set("User-Agent", "polar-client")
	req.Header.Set("Content-Type", "application/json")

	if token {
		bearer, _ := client.token()
		if bearer != "" {
			log.Debug("token: ", bearer)
			req.Header.Set("Authorization", "Bearer "+bearer)
		}
	}
	res, getErr := client.httpClient.Do(req)

	if getErr != nil {
		return nil, nil, fmt.Errorf("server req failed %w", getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	log.Debug(string(body))
	if readErr != nil {
		return nil, nil, fmt.Errorf("server req failed %w", readErr)
	}

	return body, res, nil
}

// PostReq performs a http POST request on a resource
func (client *Client) PostReq(url string, jsonStr []byte, token bool, method string) ([]byte, *http.Response, error) {
	url = client.BaseURL + url
	log.Debug("calling: ", url)
	log.Debug("payload: ", string(jsonStr))
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, nil, fmt.Errorf("server req failed %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if token {
		bearer, _ := client.token()
		if bearer != "" {
			log.Debug("token: ", bearer)
			req.Header.Set("Authorization", "Bearer "+bearer)
		}
	}
	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("server req failed %w", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, nil, fmt.Errorf("server req failed %w", readErr)
	}
	return body, res, nil
}