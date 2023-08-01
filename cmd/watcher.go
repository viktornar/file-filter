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
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func ServeWatcher(name string, command *cli.Command, arguments []string) int {
	parsed, err := command.Parse(arguments)

	if err != nil {
		return command.PrintHelp()
	}

	initLogger(parsed[2], parsed[3], name)

	logger.Printf("Starting file watcher with options: %s, %s, %s, %s\n", parsed[0], parsed[1], parsed[2], parsed[3])

	w := initWatcher(parsed[0])

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case event := <-w.Event:
				logger.Debug.Printf("Received event %v", event)
				internal.HandleWatcherEvent(event, parsed[1])
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

func initWatcher(path string) *file.Watcher {
	w := file.NewWatcher()
	w.Add(path)
	return w
}

func initLogger(logPath string, logLevel string, name string) {
	fileHandler, err := file.OpenFile(filepath.Join(logPath, fmt.Sprintf("%s_%s.log", name, strconv.Itoa(int(time.Now().Unix())))))

	if err != nil {
		panic("Unable to create a log file")
	}

	logger.SetLevel(logLevel)
	logger.SetOutput(fileHandler)
}
