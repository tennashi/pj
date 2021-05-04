package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:   "init",
	Usage:  "Initialize the project",
	Action: InitAction,
}

func InitAction(c *cli.Context) error {
	cfg, err := NewConfig(c)
	if err != nil {
		return err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	projectName := c.Args().First()

	if projectName == "" {
		projectName = filepath.Base(currentWorkingDir)
	}

	project := Project{
		Name:             projectName,
		Workspaces:       []string{currentWorkingDir},
		CurrentWorkspace: currentWorkingDir,
	}

	ret, err := client.Create(project)
	if err != nil {
		return err
	}

	var o interface{}
	outputMethod := c.String("output")
	var p Printer
	switch outputMethod {
	case "json":
		p = NewJSONPrinter(os.Stdout)
		o = ret
	default:
		p = NewTablePrinter(os.Stdout)
		o = Projects{ret}
	}

	err = p.Print(o)
	if err != nil {
		return err
	}

	return nil
}
