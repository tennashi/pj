package main

import (
	"github.com/urfave/cli/v2"
)

var ChangeCommand = &cli.Command{
	Name:   "change",
	Usage:  "Change the current project to the specified project",
	Action: ChangeAction,
}

func ChangeAction(c *cli.Context) error {
	cfg, err := NewConfig(c)
	if err != nil {
		return err
	}

	cli, err := NewClient(cfg)
	if err != nil {
		return err
	}

	projectName := c.Args().First()

	err = cli.ChangeCurrentProject(projectName)
	if err != nil {
		return err
	}

	return nil
}
