package cmd

import (
	"file-filter/pkg/cli"
	"fmt"
)

func ServeLogger(command *cli.Command, arguments []string) int {
	parsed, err := command.Parse(arguments)

	if err != nil {
		return command.PrintHelp()
	}

	fmt.Printf("%s", parsed[0])

	return cli.Success
}
