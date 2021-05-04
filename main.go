package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "HEAD"
	date    = "now"
)

func NewConfig(c *cli.Context) (*Config, error) {
	dataDir := c.String("data-dir")
	if dataDir != "" {
		return &Config{
			DataDir: dataDir,
		}, nil
	}

	dataDir = os.Getenv("XDG_DATA_HOME")
	if dataDir != "" {
		return &Config{
			DataDir: filepath.Join(dataDir, "pj"),
		}, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &Config{
		DataDir: filepath.Join(home, ".local/share/pj"),
	}, nil
}

func run() (exitCode int) {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%v version %v (rev: %v)\n", c.App.Name, c.App.Version, commit)
		fmt.Fprintf(c.App.Writer, "built at: %v\n", date)
	}

	app := cli.App{
		Name:    "pj",
		Version: version,
		Usage:   "A tool for managing units of your work",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "data-dir",
				Usage:   "specify the path to store the projects file (default: $XDG_DATA_HOME/pj)",
				EnvVars: []string{"PJ_DATA_DIR"},
			},
			&cli.StringFlag{
				Name:    "output",
				Usage:   "specify the output format (one of: json) (default: \"\")",
				Aliases: []string{"o"},
			},
		},
		Commands: []*cli.Command{
			ListCommand,
			GetCommand,
			InitCommand,
			ChangeCommand,
			CurrentCommand,
			WorkspaceCommand,
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			panic(err)
		},
	}

	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			switch c := err.(type) {
			case cli.ExitCoder:
				exitCode = c.ExitCode()
			}

			exitCode = 1
		}
	}()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
