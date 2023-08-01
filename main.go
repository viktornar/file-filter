package main

import (
	"file-filter/cmd"
	"file-filter/pkg/cli"
	"os"
)

func main() {
	os.Exit((&cli.App{
		Name:    "file-filter",
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:  "logger",
				Usage: "filter logs by given criteria",
				Arguments: []string{
					"dateFilter",
					"nameFilter",
				},
				HandleFunc: cmd.ServeLogger,
			},
			{
				Name:  "watcher",
				Usage: "watch file changes by given paths",
				Arguments: []string{
					"hotPath",
					"backupPath",
					"logPath",
					"logLevel",
				},
				HandleFunc: cmd.ServeWatcher,
			},
		},
	}).Execute())
}
