package cli

import (
	"fmt"
	"os"
)

type App struct {
	Name    string
	Version string
	Group   *Group
}

func (a *App) Execute() {
	if len(os.Args) > 1 {
		for _, command := range a.Group.Commands {
			if os.Args[1] == command.Name {
				os.Exit(command.HandleFunc(a.Group, command, os.Args[3:]))
			}

			os.Exit(a.Group.PrintHelp())
		}


		if os.Args[1] == a.Group.Name {
			if len(os.Args) > 2 {
				for _, command := range a.Group.Commands {
					if os.Args[2] == command.Name {
						os.Exit(command.HandleFunc(a.Group, command, os.Args[3:]))
					}
				}
			}

			os.Exit(a.Group.PrintHelp())
		}

		if a.Group != nil && len(a.Group.Commands) > 0 {
			group := a.Group
			command := group.Commands[0]

			os.Exit(command.HandleFunc(group, command, os.Args[1:]))
		}
	}

	os.Exit(a.PrintHelp())
}

func (a *App) PrintHelp() int {
	fmt.Printf("%s version %s\n", a.Name, a.Version)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s [flags] [arguments]\n", a.Group.Name)
	fmt.Println()

	return Success
}
