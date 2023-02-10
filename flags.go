package main

import "github.com/urfave/cli/v2"

var (
	User   string
	App    string
	unsafe bool
	Flags  []cli.Flag
)

func InitFlags() []cli.Flag {
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "user",
			Aliases:     []string{"u"},
			Value:       "",
			Usage:       "Specify the user",
			Required:    true,
			Destination: &User,
		},
		&cli.StringFlag{
			Name:        "app",
			Aliases:     []string{"a"},
			Value:       "",
			Usage:       "Specify the app",
			Destination: &App,
		},
		&cli.BoolFlag{
			Name:        "unsafe",
			Value:       false,
			Usage:       "Disable safe restrictions on password",
			Destination: &unsafe,
		},
	}
	return Flags
}
