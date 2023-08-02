package cmd

import (
	"errors"
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

		if !internal.IsValidDate(parsed[0]) {
			fmt.Printf("Date need to be in such given formats %v\n", internal.GetDateLayouts())
			return cli.Failure
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

		internal.InitLogger("debug", name)
		logPath := fmt.Sprintf("%s.log", name)

		if _, err := os.Stat(logPath); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Log file %s is not initialized yes or it is not possible to create it. Quiting...", logPath)
			return cli.Failure
		}

		w := initWatcher(logPath)

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
					fmt.Print("Closing...")
					w.Close()
					return
				case <-w.Closed:
					w.Close()
					return
				}
			}
		}()

		w.Start(time.Millisecond * 100)

		return cli.Success
	}
}
