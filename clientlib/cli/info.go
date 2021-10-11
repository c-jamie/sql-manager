package cli

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v2"
	"github.com/c-jamie/sql-manager/clientlib/app"

)

// Token returns the current application token
func Token(c *cli.Context) error {
	fmt.Println("Loading token")
	debug := ""
	if c.String("verbose") == "0" {
		debug = "info"
	} else if c.String("verbose") == "1" {
		debug = "debug"
	}
	app, err := app.New(debug)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to initialise client", err)
		return nil
	}
	dat, err := app.Token.Get()
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), err)
		return nil
	}
	result, err := json.Marshal(dat)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to load json", err)
		return nil
	}
	fmt.Println(string(pretty.Color(pretty.Pretty(result), nil)))
	return nil
}

// Version returns the current version info
func Version(c *cli.Context) error {
	debug := ""
	if c.String("verbose") == "0" {
		debug = "info"
	} else if c.String("verbose") == "1" {
		debug = "debug"
	}
	app, err := app.New(debug)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to initialise client", err)
		return nil
	}
	dat, err := app.GetVersion()
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to get version", err)
	}
	fmt.Println(cCy.Sprint("Version:"), dat)
	return nil
}

// Server returns the current server version
func Server(c *cli.Context) error {
	debug := ""
	if c.String("verbose") == "0" {
		debug = "info"
	} else if c.String("verbose") == "1" {
		debug = "debug"
	}
	app, err := app.New(debug)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to initialise client", err)
		return nil
	}
	dat, err := app.Server()
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to get version", err)
	}
	fmt.Println(cCy.Sprint("Server Info:"), dat)
	return nil
}
