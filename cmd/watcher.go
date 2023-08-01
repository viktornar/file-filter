package cmd

import (
	"file-filter/internal"
	"file-filter/pkg/cli"
	"file-filter/pkg/file"
	"file-filter/pkg/logger"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ServeWatcher(ctx *internal.Ctx) func(string, *cli.Command, []string) int {
	return func(name string, command *cli.Command, arguments []string) int {
		parsed, err := command.Parse(arguments)

		if err != nil && ctx.Watcher == (internal.Watcher{}) {
			return command.PrintHelp()
		}

		// Load initial ctx for watcher
		if ctx.Watcher == (internal.Watcher{}) {
			ctx.Watcher = internal.Watcher{
				HotPath:    parsed[0],
				BackupPath: parsed[1],
				LogLevel:   parsed[2],
			}
			internal.SaveCtx(name, ctx)
		} else {
			// Override loaded ctx with parsed arguments
			for idx, argument := range parsed {
				switch idx {
				case 0:
					ctx.Watcher.HotPath = argument
				case 1:
					ctx.Watcher.BackupPath = argument
				case 2:
					ctx.Watcher.LogLevel = argument
				}
			}
		}

		initLogger(ctx.Watcher.LogLevel, name)

		logger.Printf("Starting file watcher with options: %v\n", ctx.Watcher)

		w := initWatcher(ctx.Watcher.HotPath)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		go func() {
			for {
				select {
				case event := <-w.Event:
					logger.Debug.Printf("Received event %v", event)
					internal.HandleWatcherEvent(event, ctx.Watcher.BackupPath)
				case err := <-w.Error:
					log.Fatalln(err)
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

func initWatcher(path string) *file.Watcher {
	w := file.NewWatcher()
	w.Add(path)
	return w
}

func initLogger(logLevel string, name string) {
	fileHandler, err := file.OpenFile(fmt.Sprintf("%s.log", name))

	if err != nil {
		panic("Unable to create a log file")
	}

	logger.SetLevel(logLevel)
	logger.SetOutput(fileHandler)
}
