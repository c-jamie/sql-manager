package cli

import (
	"fmt"
	"github.com/c-jamie/sql-manager/clientlib/app"
	sqlMig "github.com/c-jamie/sql-manager/clientlib/migration"
	"github.com/c-jamie/sql-manager/clientlib/script"
	"github.com/c-jamie/sql-manager/clientlib/sql"
	"github.com/urfave/cli/v2"
)


// ScriptRegister registers a script with the platform
func ScriptRegister(c *cli.Context) error {
	debug := ""
	if c.String("verbose") == "0" {
		debug = "info"
	} else if c.String("verbose") == "1" {
		debug = "debug"
	}
	ap, nil := app.New(debug)
	err := ap.Script.Register(c.Args().First())
	if err != nil {
		fmt.Println(cRe.Sprint("Error: "), "unable to register file", err)
		return nil
	}
	fmt.Println(cGr.Sprint("Success: "), "registered")
	return nil
}


// ScriptGet returns a script from the platform
func ScriptGet(c *cli.Context) error {
	fmt.Println("Grabbing file")
	debug := ""
	if c.String("verbose") == "0" {
		debug = "info"
	} else if c.String("verbose") == "1" {
		debug = "debug"
	}
	app, nil := app.New(debug)
	sqlFile, err := app.Script.Get(c.Args().First())
	if err != nil {
		fmt.Println(cRe.Sprint("Error: "), "unable to get the script", err)
		return nil
	}
	fmt.Println(sqlFile)
	return nil
}


// ScriptGetCompile returns a script from the platform and compiles it
func ScriptGetCompile(c *cli.Context) error {
	debug := ""
	if c.String("verbose") == "0" {
		debug = "info"
	} else if c.String("verbose") == "1" {
		debug = "debug"
	}
	file := c.Args().First()
	env := c.String("env")

	app, err := app.New(debug)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to initialise client", err)
		return nil
	}
	sqlFile, err := app.Script.Get(file)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to load file", err)
		return nil
	}

	mig, err := app.Migration.GetAll(env)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to get migrations for env", err)
		return nil
	}
	migEnv := make(map[string][]*sqlMig.SQLMigrationStrategy)
	migEnv[env] = mig
	sql := sql.New(sqlFile, env, migEnv, app.Script.Get)
	sql.Compile()
	fmt.Println(sql.Parsed)
	return nil
}

// ScriptList lists all the available scripts
func ScriptList(c *cli.Context) error {
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
	project := c.Args().First()

	if project == "" {
		fmt.Println(cRe.Sprint("Error:"), "argument must be a project", err)
		return nil
	}

	files, err := app.Script.List(project)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to get tables from server", err)
		return nil
	}
	script.ProjectInfoToTable(files)
	return nil
}
