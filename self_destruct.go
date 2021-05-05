package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var SelfDestructCommand = &cli.Command{
	Name:    "self-destruct",
	Usage:   "Uninstall the pj itself",
	Aliases: []string{"self-uninstall"},
	Action:  SelfDestructAction,
}

func SelfDestructAction(c *cli.Context) error {
	cfg, err := NewConfig(c)
	if err != nil {
		return err
	}

	err = os.RemoveAll(cfg.DataDir)
	if err != nil {
		return err
	}

	cmdPath, err := os.Executable()
	if err != nil {
		return err
	}

	realCmdPath, err := filepath.EvalSymlinks(cmdPath)
	if err != nil {
		return err
	}

	err = os.Remove(realCmdPath)
	if err != nil {
		return nil
	}

	return nil
}
