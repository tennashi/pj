package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

var GetCommand = &cli.Command{
	Name:   "get",
	Usage:  "Get the project with the specified name",
	Action: GetAction,
}

func GetAction(c *cli.Context) error {
	projectName := c.Args().First()

	cfg, err := NewConfig(c)
	if err != nil {
		return err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	project, err := client.Get(projectName)
	if err != nil {
		return err
	}

	var o interface{}
	outputMethod := c.String("output")
	var p Printer
	switch outputMethod {
	case "json":
		p = NewJSONPrinter(os.Stdout)
		o = project
	default:
		p = NewTablePrinter(os.Stdout)
		o = Projects{project}
	}

	err = p.Print(o)
	if err != nil {
		return err
	}

	return nil
}
