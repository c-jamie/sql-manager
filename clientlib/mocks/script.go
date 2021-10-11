package mocks

import (
	"errors"

	scr "github.com/c-jamie/sql-manager/clientlib/script"
)

type Script struct {
}

func (app *Script) Create(name string, file string, dir string) error {
	return nil
}

func (app *Script) List(project string) ([]*scr.FileList, error) {
	return nil, nil

}

func (app *Script) Get(name string) (string, error) {
	if name == "proj1-test3-sql" {
		return `
			/*
			[sqlmbegin]
			[script]
				- description: "this is test 3"
			[dev]
				- table1: "A.B.test3"
			[sqlmend]
			*/
			select top 10 * from {{.table1}}
		`, nil
	} else {
		return "", errors.New("doesn't exist")
	}
}

func (app *Script) Register(file string) error {
	return nil

}
