package internal

import (
	"errors"
	"file-filter/pkg/file"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	deletePrefix = "delete_"
)

func HandleWatcherEvent(event file.Event, destination string) {
	fmt.Println(event)

	if event.IsDir() {
		handleDirChange(event, destination)
	} else {
		handleFileChange(event, destination)
	}
}

func handleDirChange(event file.Event, destination string) {
	// TODO: Complete logic for directory case
	switch event.Op {
	case file.Create:
		file.CreateDir(filepath.Join(destination, event.Name()))
	}
}

func handleFileChange(event file.Event, destination string) {
	// TODO: Complete logic when file is moved from one directory to another
	switch event.Op {
	case file.Create, file.Write:
		file.CopyFile(event.Path, filepath.Join(destination, event.Name()))
	case file.Rename:
		parts := strings.Split(event.Path, "/")
		oldName := event.Name()
		newName := parts[len(parts)-1]
		oldBackupPath := filepath.Join(destination, oldName)
		newBackupPath := filepath.Join(destination, newName)

		if strings.HasPrefix(newName, deletePrefix) {
			file.DeleteFile(event.Path)
			file.DeleteFile(filepath.Join(destination, oldName))
			return
		}

		if _, err := os.Stat(oldBackupPath); errors.Is(err, os.ErrNotExist) {
			file.CopyFile(event.Path, filepath.Join(destination, newName))
		}

		file.RenameFile(oldBackupPath, newBackupPath)
	case file.Remove:
		file.DeleteFile(filepath.Join(destination, event.Name()))
	}
}
