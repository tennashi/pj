package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

var WorkspaceCommand = &cli.Command{
	Name:  "workspace",
	Usage: "Subcommands for managing workspaces",
	Subcommands: []*cli.Command{
		WorkspaceAddCommand,
	},
}

var WorkspaceAddCommand = &cli.Command{
	Name:   "add",
	Usage:  "Add a workspace to the current project",
	Action: WorkspaceAddAction,
}

func WorkspaceAddAction(c *cli.Context) error {
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

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	workspacesMap := map[string]struct{}{}
	for _, workspace := range project.Workspaces {
		workspacesMap[workspace] = struct{}{}
	}

	if _, ok := workspacesMap[currentWorkingDir]; !ok {
		project.Workspaces = append(project.Workspaces, currentWorkingDir)
		project.CurrentWorkspace = currentWorkingDir
	}

	err = cli.Update(project)
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
