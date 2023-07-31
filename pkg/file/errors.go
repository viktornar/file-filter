package file

import "errors"

var (
	ErrDurationTooShort     = errors.New("file: duration is less than 1ns")
	ErrWatcherRunning       = errors.New("file: watcher is already running")
	ErrWatchedFileDeleted   = errors.New("file: watched file or folder deleted")
	ErrSkip                 = errors.New("file: skipping file")
	ErrUnableToCreateTmpDir = errors.New("file: unable to create temp directory")
)
