package cli

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/c-jamie/sql-manager/clientlib/app"
	sqlMig "github.com/c-jamie/sql-manager/clientlib/migration"
	"github.com/urfave/cli/v2"
)


// MigrationAdd add a new migration to the platform
func MigrationAdd(c *cli.Context) error {
	debug := ""
	if c.String("verbose") == "0" {
		debug = "info"
	} else if c.String("verbose") == "1" {
		debug = "debug"
	}
	file := c.Args().Get(0)

	if file == "" {
		return fmt.Errorf("file is missing")
	}
	env := c.String("env")
	table := c.String("table")

	app, err := app.New(debug)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to initialise client", err)
		return nil
	}
	err = app.Migration.Add(file, env, table)

	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to add migrate", err)
		return nil
	}
	fmt.Println(cGr.Sprint("Success:"), "migration added")
	return nil
}


// DoMigrations runs the migration for a given table and env
func DoMigrations(c *cli.Context) error {

	debug := ""
	if c.String("verbose") == "0" {
		debug = "info"
	} else if c.String("verbose") == "1" {
		debug = "debug"
	}
	table := c.Args().Get(0)

	if table == "" {
		return fmt.Errorf("table is missing")
	}
	env := c.String("env")
	connection := c.String("connection")
	driver := c.String("driver")

	app, err := app.New(debug)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to initialise client", err)
		return nil
	}

	err = app.Migration.Run(driver, connection, env, table)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to run migrations", err)
		return nil
	}
	fmt.Println(cGr.Sprint("Success:"), "migrations applied")
	return nil
}


// MigrationList list available migrations for a table and env
func MigrationList(c *cli.Context) error {
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
	table := c.Args().Get(0)
	if table == "" {
		return fmt.Errorf("table is missing")
	}
	env := c.String("env")
	mig, err := app.Migration.Get(env, table)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to initialise client", err)
		return nil
	}

	sqlMig.MigrationsToTable(*mig)
	return nil
}


// MigrationSet sets the current migration
// If you have 3 applied migrations, and you want to roll back to migration 2
// You can set migration to 2, and the next time you run migrations, you'll do a down migrtion for 3
func MigrationSet(c *cli.Context) error {
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
	env := c.Args().Get(0)
	if env == "" {
		return fmt.Errorf("env is missing")
	}
	var selectedTable string
	var tablesSelection []string

	tables, err := app.Migration.ListTables(env)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to set migration", err)
		return nil
	}
	for _, t := range tables {
		tablesSelection = append(tablesSelection, t.Table)
	}
	prompt := &survey.Select{
		Message: "Choose a migration to update",
		Options: tablesSelection,
	}
	survey.AskOne(prompt, &selectedTable)

	migrations, err := app.Migration.Get(env, selectedTable)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to set migration", err)
		return nil
	}
	var migrationIDs []string
	var selectedMigration string

	if migrations != nil {
		if migrations.MigrationsUp != nil {
			for _, m := range migrations.MigrationsUp {
				fmt.Println("adding id", m.ID)
				migrationIDs = append(migrationIDs, strconv.Itoa(m.ID))
			}
		}
		if migrations.MigrationsDown != nil {
			for _, m := range migrations.MigrationsDown {
				migrationIDs = append(migrationIDs, strconv.Itoa(m.ID))
			}
		}
	} else {
		fmt.Println(cRe.Sprint("Error:"), "unable to set migration - no migrations")
		return nil
	}

	prompt = &survey.Select{
		Message: "Choose a migration to set to",
		Options: migrationIDs,
	}
	survey.AskOne(prompt, &selectedMigration)
	selectedMigInt, err := strconv.Atoi(selectedMigration)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to set migration", err)
		return nil
	}
	err = app.Migration.Set(env, selectedTable, selectedMigInt)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to set migration", err)
		return nil
	}
	fmt.Println(cRe.Sprint("Success"), "migration updated", err)
	return nil
}

// MigrationDelete removes a migration from the platform
func MigrationDelete(c *cli.Context) error {
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
	env := c.Args().Get(0)
	if env == "" {
		return fmt.Errorf("env is missing")
	}
	var selectedTable string
	var tablesSelection []string

	tables, err := app.Migration.ListTables(env)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to delete migration", err)
		return nil
	}
	for _, t := range tables {
		tablesSelection = append(tablesSelection, t.Table)
	}
	prompt := &survey.Select{
		Message: "Choose a migration to delete",
		Options: tablesSelection,
	}
	survey.AskOne(prompt, &selectedTable)

	migrations, err := app.Migration.Get(env, selectedTable)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to delete migration", err)
		return nil
	}
	var migrationIDs []string
	var selectedMigration string

	if migrations != nil {
		if migrations.MigrationsUp != nil {
			for _, m := range migrations.MigrationsUp {
				fmt.Println("adding id", m.ID)
				migrationIDs = append(migrationIDs, strconv.Itoa(m.ID))
			}
		}
		if migrations.MigrationsDown != nil {
			for _, m := range migrations.MigrationsDown {
				migrationIDs = append(migrationIDs, strconv.Itoa(m.ID))
			}
		}
	} else {
		fmt.Println(cRe.Sprint("Error:"), "unable to set delete - no migrations")
		return nil
	}

	prompt = &survey.Select{
		Message: "Choose a migration to set to",
		Options: migrationIDs,
	}
	survey.AskOne(prompt, &selectedMigration)
	selectedMigInt, err := strconv.Atoi(selectedMigration)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to delete migration", err)
		return nil
	}
	err = app.Migration.Delete(env, selectedTable, selectedMigInt)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to delete migration", err)
		return nil
	}
	fmt.Println(cRe.Sprint("Success"), "migration deleted", err)
	return nil
}


// MigrationTableList lists all the tables we have a migration for, for a given env
func MigrationTableList(c *cli.Context) error {
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
	env := c.Args().Get(0)
	if env == "" {
		return fmt.Errorf("env is missing")
	}
	mig, err := app.Migration.ListTables(env)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to list tables", err)
		return nil
	}
	sqlMig.MigrationTablesToTable(mig)
	return nil

}