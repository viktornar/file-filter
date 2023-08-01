package cmd

import (
	"file-filter/pkg/cli"
)

func ServeLogger(name string, command *cli.Command, arguments []string) int {
	_, err := command.Parse(arguments)

	if err != nil {
		return command.PrintHelp()
	}

	return cli.Success
}
