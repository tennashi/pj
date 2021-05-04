# pj [![Actions status](https://github.com/tennashi/pj/workflows/test/badge.svg)](https://github.com/tennashi/pj/actions)
A tool for managing units of your work

## Usage
```
NAME:
   pj - A tool for managing units of your work

USAGE:
   pj [global options] command [command options] [arguments...]

COMMANDS:
   list       List projects
   get        Get the project with the specified name
   init       Initialize the project
   change     Change the current project to the specified project
   current    Get the current project
   workspace  Subcommands for managing workspaces
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --data-dir value          specify the path to store the projects file (default: $XDG_DATA_HOME/pj) [$PJ_DATA_DIR]
   --output value, -o value  specify the output format (one of: json) (default: "")
   --help, -h                show help (default: false)
```

## License
MIT
