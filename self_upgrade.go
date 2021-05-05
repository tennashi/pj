package main

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli/v2"
)

var SelfUpgradeCommand = &cli.Command{
	Name:   "self-upgrade",
	Usage:  "Upgrade the pj itself",
	Action: SelfUpgradeAction,
}

func SelfUpgradeAction(c *cli.Context) error {
	v, err := semver.Parse(version)
	if err != nil {
		if version == "dev" {
			fmt.Println("this is development version, so didn't upgrade it.")
			return nil
		}
		return err
	}

	latest, err := selfupdate.UpdateSelf(v, "tennashi/pj")
	if err != nil {
		return err
	}

	if latest.Version.Equals(v) {
		fmt.Println("up to date")
		return nil
	}

	fmt.Println("successfully updated to version", latest.Version)
	fmt.Println("release note:")
	fmt.Print(latest.ReleaseNotes)

	return nil
}
