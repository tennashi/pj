package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

var CurrentCommand = &cli.Command{
	Name:   "current",
	Usage:  "Get the current project",
	Action: CurrentAction,
}

func CurrentAction(c *cli.Context) error {
	cfg, err := NewConfig(c)
	if err != nil {
		return err
	}

	cli, err := NewClient(cfg)
	if err != nil {
		return err
	}

	project, err := cli.GetCurrentProject()
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
