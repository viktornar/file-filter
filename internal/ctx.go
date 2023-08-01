package internal

type Ctx struct {
	Watcher
	Filter
}

type Watcher struct {
	HotPath    string
	BackupPath string
	LogLevel   string
	LogPath    string
}

type Filter struct {
	Name string
	Date string
}
