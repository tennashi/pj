package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

var ListCommand = &cli.Command{
	Name:   "list",
	Usage:  "List projects",
	Action: ListAction,
}

func ListAction(c *cli.Context) error {
	cfg, err := NewConfig(c)
	if err != nil {
		return err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	projects, err := client.List()
	if err != nil {
		return err
	}

	var o interface{}
	outputMethod := c.String("output")
	var p Printer
	switch outputMethod {
	case "json":
		p = NewJSONPrinter(os.Stdout)
		o = projects
	default:
		p = NewTablePrinter(os.Stdout)

		o = projects
	}

	err = p.Print(o)
	if err != nil {
		return err
	}

	return nil
}
