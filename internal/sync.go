package internal

import (
	"errors"
	"file-filter/pkg/file"
	"os"
	"path/filepath"
	"strings"
)

const (
	deletePrefix    = "delete_"
	backupExtension = ".bak"
)

func HandleWatcherEvent(event file.Event, destination string) {
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
		file.CopyFile(event.Path, filepath.Join(destination, event.Name()+backupExtension))
	case file.Rename:
		_, newName := filepath.Split(event.Path)

		oldName := event.Name()
		oldBackupPath := filepath.Join(destination, oldName+backupExtension)
		newBackupPath := filepath.Join(destination, newName+backupExtension)

		if strings.HasPrefix(newName, deletePrefix) {
			file.DeleteFile(event.Path)
			file.DeleteFile(oldBackupPath)
			return
		}

		if _, err := os.Stat(oldBackupPath); errors.Is(err, os.ErrNotExist) {
			file.CopyFile(event.Path, newBackupPath)
			return
		}

		file.RenameFile(oldBackupPath, newBackupPath)
	case file.Remove:
		file.DeleteFile(filepath.Join(destination, event.Name()+backupExtension))
	}
}
