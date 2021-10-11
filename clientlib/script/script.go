package script

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/c-jamie/sql-manager/clientlib/request"
	"github.com/jedib0t/go-pretty/table"
	"github.com/tidwall/gjson"
)

// FileList represents a given SQL script
type FileList struct {
	ProjectID    int64  `json:"project_id"`
	ProjectName  string `json:"project_name"`
	FileID       int64  `json:"file_id"`
	FileLocation string `json:"file_location"`
	SnippitName  string `json:"snippit_name"`
	SnippitID    int64  `json:"snippit_id"`
}


// Script represents the interface for interacting with SQL Scripts
type Script interface {
	// List returns all scripts for a project
	List(project string) ([]*FileList, error)
	// Get returns a script for a given name
	Get(name string) (string, error)
	// Register registers a script
	Register(file string) error
}

type script struct {
	BaseURL string
}

const (
	File = "files"
)


// New creates a new script object
func New(url string) Script {
	scr := &script{BaseURL: url}
	return scr
}

func (app *script) makeRequest(url string, payload []byte, how string) ([]byte, error) {
	client := request.Client{BaseURL: app.BaseURL}
	var body []byte
	var resp *http.Response
	var err error
	if how == http.MethodPost {
		body, resp, err = client.PostReq(url, payload, true, how)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusCreated {
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

func (app *script) Get(name string) (string, error) {
	url := "/" + File + "?name=" + name
	body, err := app.makeRequest(url, nil, http.MethodGet)
	if err != nil {
		return "", err
	}
	fileRe := gjson.Get(string(body), File+".script").String()
	return fileRe, nil
}

func (app *script) List(project string) ([]*FileList, error) {
	unmarshal := func(body []byte) ([]*FileList, error) {
		var projInf []*FileList
		jsonErr := json.Unmarshal(body, &projInf)
		if jsonErr != nil {
			return nil, jsonErr
		} else {
			return projInf, nil
		}
	}
	url := "/" + File + "/list?project=" + project
	body, err := app.makeRequest(url, make([]byte, 0), http.MethodGet)
	if err != nil {
		return nil, err
	}
	files := gjson.Get(string(body), "files").String()
	proInf, err := unmarshal([]byte(files))
	if err != nil {
		return nil, fmt.Errorf("unable to list tables %w", err)
	} else {
		return proInf, nil
	}
}

func (app *script) Register(file string) error {
	url := "/" + File
	var jsonStr = []byte(fmt.Sprintf(`{"file":"%s", "project":"%s"}`, file, filepath.Dir(file)))
	_, err := app.makeRequest(url, jsonStr, http.MethodPost)
	if err != nil {
		return err
	}
	return nil
}

func ProjectInfoToTable(proj []*FileList) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "File ID", "File Location", "File Name"})
	for i, k := range proj {
		t.AppendRow(table.Row{i, k.FileID, k.FileLocation, k.SnippitName})
	}
	t.Render()

}
