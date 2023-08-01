package cli

import (
	"fmt"
	"os"
	"strings"
)

type App struct {
	Name     string
	Version  string
	Commands []*Command
}

func (a *App) Execute() int {
	if len(os.Args) > 1 {
		for _, command := range a.Commands {
			if os.Args[1] == command.Name {
				return command.HandleFunc(a.Name, command, os.Args[2:])
			}
		}
	}

	return a.PrintHelp()
}

func (a *App) PrintHelp() int {
	fmt.Printf("%s version %s\n", a.Name, a.Version)
	fmt.Println()
	fmt.Println("Usage:")

	commands := []string{}

	for _, command := range a.Commands {
		commands = append(commands, command.Name)
	}

	fmt.Printf("  %s <%s> <arguments>\n", a.Name, strings.Join(commands, "|"))
	fmt.Println()

	return Success
}
