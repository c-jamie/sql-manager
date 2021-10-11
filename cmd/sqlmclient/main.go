package main

import (
	"fmt"
	"log"
	"os"

	smcli "github.com/c-jamie/sql-manager/clientlib/cli"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "sqlmclient",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "json",
				Value:   false,
				Aliases: []string{"j"},
				Usage:   "return all output as json",
			},
			&cli.IntFlag{
				Name:    "verbose",
				Value:   0,
				Aliases: []string{"v"},
				Usage:   "0 no logging, 1 for info, 2 for debug",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "login",
				Usage:  "login to client",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "password",
						Aliases:  []string{"p"},
						Required: false,
						Usage:    "your password",
						Value:    "none",
					},
				},
				Action: smcli.Login,
			},
			{
				Name:  "register",
				Usage: "Register with SQL Manager",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "team",
						Aliases:  []string{"n"},
						Required: false,
						Usage:    "the team name",
						Value:    "none",
					},
					&cli.StringFlag{
						Name:     "email",
						Aliases:  []string{"e"},
						Required: false,
						Usage:    "the email address of the user",
						Value:    "none",
					},
					&cli.StringFlag{
						Name:     "password",
						Aliases:  []string{"p"},
						Required: false,
						Usage:    "the user password",
						Value:    "none",
					},
				},
				Action: smcli.Register,
			},
			{
				Name:  "script",
				Usage: "work with your SQL",
				Subcommands: []*cli.Command{
					{
						Name:    "register",
						Aliases: []string{"r"},
						Usage:   "register a file with the service",
						Action:  smcli.ScriptRegister,
					},
					{
						Name:    "remove",
						Aliases: []string{"rm"},
						Usage:   "deregister a file from the application",
						Action: func(c *cli.Context) error {
							fmt.Println("completed task: ", c.Args().First())
							return nil
						},
					},
					{
						Name:    "get",
						Aliases: []string{"g"},
						Usage:   "return a file fully parsed SQL script",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "env",
								Required: false,
								Aliases:  []string{"e"},
								Value:    "none",
								Usage:    "the environment context used to parse the file",
							},
						},
						Action: smcli.ScriptGet,
					},
					{
						Name:    "list",
						Aliases: []string{"ls"},
						Usage:   "get all available SQL scripts",
						Action:  smcli.ScriptList,
					},
					{
						Name:    "get-compile",
						Aliases: []string{"gc"},
						Usage:   "compile a script, useful for testing",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "env",
								Required: true,
								Aliases:  []string{"e"},
								Value:    "none",
								Usage:    "environment",
							},
						},
						Action: smcli.ScriptGetCompile,
					},
				},
			},
			{
				Name:  "migration",
				Usage: "add SQL migrations",
				Subcommands: []*cli.Command{
					{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "create a file with the service",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "table",
								Aliases:  []string{"t"},
								Required: true,
								Usage:    "the table name to associate",
							},
							&cli.StringFlag{
								Name:     "env",
								Aliases:  []string{"e"},
								Required: true,
								Usage:    "the env",
							},
						},
						Action: smcli.MigrationAdd,
					},
					{
						Name:    "list",
						Aliases: []string{"ls"},
						Usage:   "list migrations",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "env",
								Aliases:  []string{"e"},
								Required: true,
								Usage:    "the env",
							},
						},
						Action: smcli.MigrationList,
					},
					{
						Name:    "set",
						Aliases: []string{"s"},
						Usage:   "list migrations",
						Action:  smcli.MigrationSet,
					},
					{
						Name:    "delete",
						Aliases: []string{"d"},
						Usage:   "delete migration",
						Action:  smcli.MigrationDelete,
					},
					{
						Name:    "list-tables",
						Aliases: []string{"lt"},
						Usage:   "list migration tables",
						Action:  smcli.MigrationTableList,
					},
					{
						Name:    "run",
						Aliases: []string{"r"},
						Usage:   "list migration tables",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "env",
								Aliases:  []string{"e"},
								Required: true,
								Usage:    "the env",
							},
							&cli.StringFlag{
								Name:     "connection",
								Aliases:  []string{"c"},
								Required: true,
								Usage:    "the connection string",
							},
							&cli.StringFlag{
								Name:     "driver",
								Aliases:  []string{"d"},
								Required: true,
								Usage:    "the driver",
							},
						},
						Action: smcli.DoMigrations,
					},
				},
			},
			{
				Name:  "info",
				Usage: "options for info",
				Subcommands: []*cli.Command{
					{
						Name:  "env",
						Usage: "get env info",
						Action: func(c *cli.Context) error {
							server := os.Getenv("SQL_MNGR_SERVER")
							central := os.Getenv("SQL_MNGR_CENTRAL")
							home := os.Getenv("SQL_MNGR_HOME")
							mode := os.Getenv("SQL_MNGR_MODE")
							fmt.Println("Server: ", server)
							fmt.Println("Central: ", central)
							fmt.Println("Home: ", home)
							fmt.Println("Mode: ", mode)
							return nil
						},
					},
					{
						Name:   "token",
						Usage:  "get client token",
						Action: smcli.Token,
					},
					{
						Name:   "user",
						Usage:  "get user",
						Action: smcli.User,
					},
					{
						Name:   "version",
						Usage:  "get version info",
						Action: smcli.Version,
					},
					{
						Name:   "server",
						Usage:  "get server info",
						Action: smcli.Server,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.
			Fatal(err)
	}
}
