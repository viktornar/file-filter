package cli

import (
	"fmt"
	"os"
	"strings"
)

type App struct {
	Name       string
	Version    string
	TryDefault bool
	Groups     []*Group
}

func (a *App) Execute() {
	if len(os.Args) > 1 {
		for _, group := range a.Groups {
			if os.Args[1] == group.Name {
				if len(os.Args) > 2 {
					for _, command := range group.Commands {
						if os.Args[2] == command.Name {
							os.Exit(command.HandleFunc(group, command, os.Args[3:]))
						}
					}
				}

				os.Exit(group.PrintHelp())
			}
		}

		if a.TryDefault && len(a.Groups) > 0 && len(a.Groups[0].Commands) > 0 {
			group := a.Groups[0]
			command := group.Commands[0]

			os.Exit(command.HandleFunc(group, command, os.Args[1:]))
		}
	}

	os.Exit(a.PrintHelp())
}

func (a *App) PrintHelp() int {
	fmt.Println(a.Name + " version " + a.Version)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  group command [flags] [arguments]")
	fmt.Println()
	fmt.Println("Available groups:")

	max := 0

	for _, group := range a.Groups {
		if max < len(group.Name) {
			max = len(group.Name)
		}
	}

	max += 2

	for _, group := range a.Groups {
		fmt.Println("  " + group.Name + strings.Repeat(" ", max-len(group.Name)) + group.Usage)
	}

	return Success
}
