package cmd

import (
	"file-filter/internal"
	"file-filter/pkg/cli"
	"file-filter/pkg/file"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ServeLogger(ctx *internal.Ctx) func(string, *cli.Command, []string) int {
	return func(name string, command *cli.Command, arguments []string) int {
		parsed, err := command.Parse(arguments)

		if err != nil && ctx.Filter == (internal.Filter{}) {
			return command.PrintHelp()
		}

		if ctx.Filter == (internal.Filter{}) {
			ctx.Filter = internal.Filter{
				Date: parsed[0],
				Name: parsed[1],
			}
			internal.SaveCtx(name, ctx)
		} else {
			for idx, argument := range parsed {
				switch idx {
				case 0:
					ctx.Filter.Date = argument
				case 1:
					ctx.Filter.Name = argument
				}
			}
		}

		fileHandler, err := file.OpenFile(fmt.Sprintf("%s.log", name))

		if err != nil {
			fmt.Println(file.ErrUnableToReadFile)
			return cli.Failure
		}

		defer fileHandler.Close()

		file.FileScanner(fileHandler, func(line string) {
			internal.PrintFilteredLine(line, &ctx.Filter)
			fmt.Println()
		})

		w := initWatcher(fmt.Sprintf("%s.log", name))

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		go func() {
			for {
				select {
				case <-w.Event:
					file.ReadLastLine(fileHandler, func(line string) {
						internal.PrintFilteredLine(line, &ctx.Filter)
					})
				case <-w.Error:
					w.Close()
					return
				case sig := <-interrupt:
					fmt.Printf("Watcher received signal %v\n", sig)
					fmt.Println("Closing...")
					w.Close()
					return
				case <-w.Closed:
					return
				}
			}
		}()

		w.Start(time.Millisecond * 100)

		return cli.Success
	}
}
