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

func (a *App) Execute() {
	if len(os.Args) > 1 {
		for _, command := range a.Commands {
			if os.Args[1] == command.Name {
				os.Exit(command.HandleFunc(command, os.Args[2:]))
			}
		}
	}

	os.Exit(a.PrintHelp())
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
