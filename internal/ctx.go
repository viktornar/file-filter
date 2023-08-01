package internal

import (
	"encoding/json"
	"file-filter/pkg/file"
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
	data, err := file.ReadFile(fmt.Sprintf("%s.json", name))

	if err != nil {
		return ctx
	}

	if err := json.Unmarshal(data, &ctx); err != nil {
		return ctx
	}

	return ctx
}

func SaveCtx(name string, ctx *Ctx) error {
	data, err := json.Marshal(ctx)

	if err != nil {
		return err
	}

	return file.WriteFile(fmt.Sprintf("%s.json", name), data)
}
