package main

import (
	"file-filter/cmd"
	"file-filter/pkg/cli"
)

func main() {
	(&cli.App{
		Name:    "FileFilter",
		Version: "0.0.1",
		Group: &cli.Group{
			Name:  "start",
			Usage: "starts file-filter commands",
			Commands: []*cli.Command{
				{
					Name:  "logger",
					Usage: "filter logs by given criteria",
					Arguments: []string{
						"dateFilter",
						"regexFilter",
					},
					HandleFunc: cmd.ServeLogger,
				},
				{
					Name:  "watcher",
					Usage: "watch file changes by given paths",
					Arguments: []string{
						"hotPath",
						"backupPath",
					},
					HandleFunc: cmd.ServeWatcher,
				},
			},
		},
	}).Execute()
}
