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
				Name:  "logger",
				Usage: "Logger commands",
				Commands: []*cli.Command{
					{
						Name:  "--filter",
						Usage: "filter by given criteria.",
						Arguments: []string{
							"date",
						},
						HandleFunc: cmd.ServeLogger,
					},
				},
			},
			{
				Name:  "watcher",
				Usage: "File watcher commands",
				Commands: []*cli.Command{
					{
						Name:       "--hot",
						Usage:      "path to the hot folder",
						HandleFunc: cmd.ServeWatcher,
					},
					{
						Name:       "--backup",
						Usage:      "path to the backup folder",
						HandleFunc: cmd.ServeWatcher,
					},
				},
			},
		},
	}).Execute()
}
