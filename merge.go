package main

import "github.com/urfave/cli/v2"

var MergeCommand = &cli.Command{
	Name:   "merge",
	Usage:  "Merge the specified project into the current project",
	Action: MergeAction,
}

func MergeAction(c *cli.Context) error {
	cfg, err := NewConfig(c)
	if err != nil {
		return err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	curProj, err := client.GetCurrentProject()
	if err != nil {
		return err
	}

	projectName := c.Args().First()
	proj, err := client.Get(projectName)
	if err != nil {
		return err
	}

	workspaceMap := make(map[string]struct{})
	for _, w := range curProj.Workspaces {
		workspaceMap[w] = struct{}{}
	}

	for _, w := range proj.Workspaces {
		workspaceMap[w] = struct{}{}
	}

	mergedWorkspaces := make([]string, 0, len(workspaceMap))
	for w := range workspaceMap {
		mergedWorkspaces = append(mergedWorkspaces, w)
	}

	curProj.Workspaces = mergedWorkspaces

	err = client.Update(curProj)
	if err != nil {
		return err
	}

	err = client.Remove(projectName)
	if err != nil {
		return err
	}

	return nil
}
