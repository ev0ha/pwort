package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ev0/pwort/actions"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func main() {
	db, err := sql.Open("sqlite3", "./passwords.db")
	if err != nil {
		log.Fatal("Could not open database.")
	}
	defer db.Close()

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, pw TEXT, UNIQUE (name))")
	if err != nil {
		log.Fatal("Could not execute sql statement for users")
	}
	defer stmt.Close()

	stmt.Exec()

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS apps (id INTEGER PRIMARY KEY, name TEXT, pw TEXT, user TEXT, FOREIGN KEY (user) REFERENCES users (name))")
	if err != nil {
		log.Fatal("Could not execute sql statement for apps")
	}
	stmt.Exec()

	flags := InitFlags()

	app := &cli.App{
		Name:  "pwort",
		Usage: "Simple cli password manager",
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"cr"},
				Usage:   "Create a new user with a master password or a new app with a password.",
				Flags:   flags,
				Action: func(cCtx *cli.Context) error {
					if err := actions.Create(&User, &App, &unsafe); err != nil {
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:    "update",
				Aliases: []string{"up"},
				Usage:   "Update user or app.",
				Flags:   flags,
				Action: func(cCtx *cli.Context) error {
					if err := actions.Update(&User, &App, &unsafe); err != nil {
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:    "show",
				Aliases: []string{"sh"},
				Usage:   "Displays the password of specified app for specified user or checks whether user exists if only user is specified.",
				Flags:   flags,
				Action: func(cCtx *cli.Context) error {
					if err := actions.Show(&User, &App); err != nil {
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"dl"},
				Usage:   "Delete an user or just a specific app for an user.",
				Flags:   flags,
				Action: func(cCtx *cli.Context) error {
					if err := actions.Delete(&User, &App); err != nil {
						os.Exit(1)
					}
					return nil
				},
			},
		},
		Action: func(cCtx *cli.Context) error {
			cli.ShowAppHelpAndExit(cCtx, 0)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
