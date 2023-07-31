package file

import (
	"os"
	"time"
)

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     interface{}
	dir     bool
}

func (fs *fileInfo) IsDir() bool {
	return fs.dir
}

func (fs *fileInfo) ModTime() time.Time {
	return fs.modTime
}

func (fs *fileInfo) Mode() os.FileMode {
	return fs.mode
}

func (fs *fileInfo) Name() string {
	return fs.name
}

func (fs *fileInfo) Size() int64 {
	return fs.size
}

func (fs *fileInfo) Sys() interface{} {
	return fs.sys
}