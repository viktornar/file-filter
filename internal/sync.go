package internal

import (
	"errors"
	"file-filter/pkg/file"
	"file-filter/pkg/logger"
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
		path := filepath.Join(destination, event.Name())
		logger.Printf("Directory event was %s. Creating directory %s in %s", file.Operations[event.Op], event.Name(), path)
		file.CreateDir(path)
	}
}

func handleFileChange(event file.Event, destination string) {
	// TODO: Complete logic when file is moved from one directory to another
	switch event.Op {
	case file.Create, file.Write:
		createWritePath := filepath.Join(destination, event.Name()+backupExtension)
		logger.Printf("File event was %s. Copying file %s from %s to %s", file.Operations[event.Op], event.Name(), event.Path, createWritePath)
		file.CopyFile(event.Path, createWritePath)
	case file.Rename:
		_, newName := filepath.Split(event.Path)

		oldName := event.Name()
		oldBackupPath := filepath.Join(destination, oldName+backupExtension)
		newBackupPath := filepath.Join(destination, newName+backupExtension)

		if strings.HasPrefix(newName, deletePrefix) {
			logger.Printf("File event was %s. File has prefix %s. Deleting files %s, %s", file.Operations[event.Op], deletePrefix, event.Path, oldBackupPath)
			file.DeleteFile(event.Path)
			file.DeleteFile(oldBackupPath)
			return
		}

		if _, err := os.Stat(oldBackupPath); errors.Is(err, os.ErrNotExist) {
			logger.Printf("File event was %s. Not found in backup. Copying file %s from %s to %s", file.Operations[event.Op], event.Name(), event.Path, oldBackupPath, newBackupPath)
			file.CopyFile(event.Path, newBackupPath)
			return
		}

		logger.Printf("File event was %s. Renaming file from %s to %s", file.Operations[event.Op], oldBackupPath, newBackupPath)
		file.RenameFile(oldBackupPath, newBackupPath)
	case file.Remove:
		backupFilePath := filepath.Join(destination, event.Name()+backupExtension)
		logger.Printf("File event was %s. Removing files %s, %s", file.Operations[event.Op], event.Path, backupFilePath)
		file.DeleteFile(backupFilePath)
	}
}
