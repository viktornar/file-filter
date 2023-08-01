package cli

import (
	"os"
	"testing"
)

func setupCliTest(t testing.TB) *App {
	app := &App{
		Name:    "file-filter",
		Version: "0.0.1",
		Commands: []*Command{
			{
				Name:  "logger",
				Usage: "filter logs by given criteria",
				Arguments: []string{
					"dateFilter",
					"regexFilter",
				},
				HandleFunc: handleFuncMock,
			},
			{
				Name:  "watcher",
				Usage: "watch file changes by given paths",
				Arguments: []string{
					"hotPath",
					"backupPath",
				},
				HandleFunc: handleFuncMock,
			},
		},
	}

	return app
}

func TestWatcherCli(t *testing.T) {
	app := setupCliTest(t)
	os.Args = []string{"file-filter", "watcher", "./test/hot", "./test/backup"}

	code := app.Execute()

	if code == Failure {
		t.Fatalf("expected to correctly parse cli options")
	}
}

func TestLoggerCli(t *testing.T) {
	app := setupCliTest(t)
	os.Args = []string{"file-filter", "logger", "2023-08-01", "test"}

	code := app.Execute()

	if code == Failure {
		t.Fatalf("expected to correctly parse cli options")
	}
}

func handleFuncMock(command *Command, arguments []string) int {
	_, err := command.Parse(arguments)

	if err != nil {
		return command.PrintHelp()
	}

	return Success
}
