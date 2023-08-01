package main

import (
	"file-filter/cmd"
	"file-filter/internal"
	"file-filter/pkg/cli"
	"os"
)

func main() {
	name := "file-filter"
	ctx := internal.LoadCtx(name)

	os.Exit((&cli.App{
		Name:    name,
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:  "logger",
				Usage: "filter logs by given criteria",
				Arguments: []string{
					"dateFilter",
					"nameFilter",
				},
				HandleFunc: cmd.ServeLogger(&ctx),
			},
			{
				Name:  "watcher",
				Usage: "watch file changes by given paths",
				Arguments: []string{
					"hotPath",
					"backupPath",
					"logLevel",
				},
				HandleFunc: cmd.ServeWatcher(&ctx),
			},
		},
	}).Execute())
}
