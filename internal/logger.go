package internal

import (
	"file-filter/pkg/file"
	"file-filter/pkg/logger"
	"fmt"
)

func InitLogger(logLevel string, name string) {
	fileHandler, err := file.OpenFile(fmt.Sprintf("%s.log", name))

	if err != nil {
		panic("Unable to create a log file")
	}

	logger.SetLevel(logLevel)
	logger.SetOutput(fileHandler)
}