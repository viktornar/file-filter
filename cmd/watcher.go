package cmd

import "file-filter/pkg/cli"

func ServeWatcher(group *cli.Group, command *cli.Command, arguments []string) int {
	println("Starting filter with options...")

	return cli.Success
}
