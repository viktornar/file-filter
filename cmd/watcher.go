package cmd

import (
	"file-filter/internal"
	"file-filter/pkg/cli"
	"file-filter/pkg/file"
	"fmt"
	"log"
	"time"
)

func ServeWatcher(group *cli.Group, command *cli.Command, arguments []string) int {
	parsed, err := command.Parse(arguments)

	if err != nil {
		return command.PrintHelp(group)
	}

	fmt.Printf("Starting file watcher with options: %s, %s\n", parsed[0], parsed[1])

	w := file.NewWatcher()

	w.Add(parsed[0])

	go func() {
		for {
			select {
			case event := <-w.Event:
				internal.HandleWatcherEvent(event, parsed[1])
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	w.Start(time.Millisecond * 100)

	return cli.Success
}

