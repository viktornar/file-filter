package main

import (
	"file-filter/cmd"
	"file-filter/pkg/cli"
)

func main() {
	(&cli.App{
		Name:    "FileFilter",
		Version: "0.0.1",
		Groups: []*cli.Group{
			{
				Name:  "start",
				Usage: "File watcher commands",
				Commands: []*cli.Command{
					{
						Name:  "logger",
						Usage: "path to the hot folder",
						Arguments: []string{
							"date",
							"regex",
						},
						HandleFunc: cmd.ServeLogger,
					},
					{
						Name:  "watcher",
						Usage: "path to the hot folder",
						Arguments: []string{
							"hot",
							"backup",
						},
						HandleFunc: cmd.ServeWatcher,
					},
				},
			},
		},
	}).Execute()
}
