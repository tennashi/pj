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
		WorkspaceListCommand,
		WorkspaceChangeCommand,
	},
}

var WorkspaceListCommand = &cli.Command{
	Name:   "list",
	Usage:  "List workspaces for the current project",
	Action: WorkspaceListAction,
}

func WorkspaceListAction(c *cli.Context) error {
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

	wss := make(Workspaces, 0, len(project.Workspaces))
	for _, ws := range project.Workspaces {
		wss = append(wss, &Workspace{
			Path: ws,
		})
	}

	var o interface{}
	outputMethod := c.String("output")
	var p Printer
	switch outputMethod {
	case "json":
		p = NewJSONPrinter(os.Stdout)
		o = wss
	default:
		p = NewTablePrinter(os.Stdout)
		o = wss
	}

	err = p.Print(o)
	if err != nil {
		return err
	}
	return nil
}

type Workspace struct {
	Path string `json:"path"`
}

func (w *Workspace) Summary() []string {
	if w == nil {
		return nil
	}

	return []string{
		w.Path,
	}
}

type Workspaces []*Workspace

func (w Workspaces) Header() []string {
	return []string{"PATH"}
}

func (w Workspaces) Summaries() []Summarable {
	if w == nil {
		return nil
	}

	ret := make([]Summarable, 0, len(w))
	for _, ws := range w {
		ret = append(ret, ws)
	}

	return ret
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

var WorkspaceChangeCommand = &cli.Command{
	Name:   "change",
	Usage:  "Change the current workspace to the specified workspace",
	Action: WorkspaceChangeAction,
}

func WorkspaceChangeAction(c *cli.Context) error {
	cfg, err := NewConfig(c)
	if err != nil {
		return err
	}

	cli, err := NewClient(cfg)
	if err != nil {
		return err
	}

	workspaceName := c.Args().First()

	err = cli.ChangeCurrentWorkspace(workspaceName)
	if err != nil {
		return err
	}

	return nil
}
