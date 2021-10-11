package cli

import (
	"fmt"
	"encoding/json"
	"github.com/AlecAivazis/survey/v2"
	"github.com/c-jamie/sql-manager/clientlib/app"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v2"
)


// User returns the acccount information for a given user
func User(c *cli.Context) error {
	fmt.Println("Grabbing account info")
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
	result, err := app.Account.User()
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), err)
		return nil
	}
	resultPretty, err := json.Marshal(result)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to convert server response to output")
		return nil
	}
	fmt.Println(string(pretty.Color(pretty.Pretty(resultPretty), nil)))
	return nil
}


// Register registers the user with the system
func Register(c *cli.Context) error {
	debug := ""
	if c.String("verbose") == "0" {
		debug = "info"
	} else if c.String("verbose") == "1" {
		debug = "debug"
	}
	team := c.String("team")
	email := c.String("email")
	password := c.String("password")

	app, err := app.New(debug)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to initialise client", err)
		return nil
	}
	if email == "none" {
		promptEmail := &survey.Input{
			Message: "please enter the email of the user to add",
		}
		survey.AskOne(promptEmail, &email)
	}
	if team == "none" {
		promptTeam := &survey.Input{
			Message: "please enter the team of the user to add",
		}
		survey.AskOne(promptTeam, &team)
	}
	if password == "none" {
		promptPw := &survey.Input{
			Message: "please enter the password",
		}
		survey.AskOne(promptPw, &password)
	}
	user, err := app.Account.Register(email, team, password)

	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	err = app.Token.AddUser(user)

	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}
	fmt.Println("The user has been", cGr.Sprint("successfully"), "added")
	return nil
}


// Login logs the user in
func Login(c *cli.Context) error {
	debug := ""
	if c.String("verbose") == "0" {
		debug = "info"
	} else if c.String("verbose") == "1" {
		debug = "debug"
	}
	password := c.String("password")

	app, err := app.New(debug)
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to initialise client", err)
		return nil
	}
	if password == "none" {
		promptPw := &survey.Input{
			Message: "please enter the password",
		}
		survey.AskOne(promptPw, &password)
	}

	token, err := app.Token.Get()
	if err != nil {
		fmt.Println(cRe.Sprint("Error:"), "unable to get token", err)
		return nil
	}

	authToken, err := app.Account.Login(token.User.Email, password)

	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	err = app.Token.AddToken(authToken)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	fmt.Println("The user has been", cGr.Sprint("successfully"), "added")
	return nil
}