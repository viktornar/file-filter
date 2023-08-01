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


func TestWatcherCliWithAllArguments(t *testing.T) {
	app := setupCliTest(t)
	os.Args = []string{"file-filter", "watcher", "./test/hot", "./test/backup"}

	code := app.Execute()

	if code == Failure {
		t.Fatalf("expected to correctly parse cli options")
	}
}

func TestLoggerCliWithAllArguments(t *testing.T) {
	app := setupCliTest(t)
	os.Args = []string{"file-filter", "logger", "2023-08-01", "test"}

	code := app.Execute()

	if code == Failure {
		t.Fatalf("expected to fail parse cli options")
	}
}

func TestWatcherCliWithMissingArguments(t *testing.T) {
	app := setupCliTest(t)
	os.Args = []string{"file-filter", "watcher", "./test/hot"}

	code := app.Execute()

	if code == Success {
		t.Fatalf("expected to failure due to missing argumets")
	}
}

func TestLoggerCliWithMissingArguments(t *testing.T) {
	app := setupCliTest(t)
	os.Args = []string{"file-filter", "logger", "2023-08-01"}

	code := app.Execute()

	if code == Success {
		t.Fatalf("expected to failure due to missing argumets")
	}
}

func handleFuncMock(command *Command, arguments []string) int {
	_, err := command.Parse(arguments)

	if err != nil {
		return Failure
	}

	return Success
}
