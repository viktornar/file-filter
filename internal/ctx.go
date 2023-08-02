package internal

import (
	"encoding/json"
	"file-filter/pkg/file"
	"file-filter/pkg/logger"
	"fmt"
)

type Ctx struct {
	Watcher `json:"watcher,omitempty"`
	Filter  `json:"filter,omitempty"`
}

type Watcher struct {
	HotPath    string `json:"hotPath,omitempty"`
	BackupPath string `json:"backupPath,omitempty"`
	LogLevel   string `json:"logLevel,omitempty"`
}

type Filter struct {
	Name string `json:"name,omitempty"`
	Date string `json:"date,omitempty"`
}

func LoadCtx(name string) Ctx {
	ctx := Ctx{}
	contextPath := fmt.Sprintf("%s.json", name)
	logger.Debug.Printf("Loading watcher context %s", contextPath)
	data, err := file.ReadFile(contextPath)

	if err != nil {
		logger.Error.Print(err)
		return ctx
	}

	if err := json.Unmarshal(data, &ctx); err != nil {
		logger.Error.Print(err)
		return ctx
	}

	return ctx
}

func SaveCtx(name string, ctx *Ctx) error {
	data, err := json.Marshal(ctx)

	logger.Debug.Print("Trying to save application context")

	if err != nil {
		logger.Error.Print(err)
		return err
	}

	return file.WriteFile(fmt.Sprintf("%s.json", name), data)
}
