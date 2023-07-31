package file

import (
	"os"
	"time"
)

type FileInfoMock struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     interface{}
	dir     bool
}

func NewFileInfoMock(name string) *FileInfoMock {
	modTime := time.Now()

	return &FileInfoMock{
		name:    name,
		size:    1,
		mode:    os.ModeDir,
		modTime: modTime,
		sys:     nil,
		dir:     true,
	}
}

func (fs *FileInfoMock) IsDir() bool {
	return fs.dir
}

func (fs *FileInfoMock) ModTime() time.Time {
	return fs.modTime
}

func (fs *FileInfoMock) Mode() os.FileMode {
	return fs.mode
}

func (fs *FileInfoMock) Name() string {
	return fs.name
}

func (fs *FileInfoMock) Size() int64 {
	return fs.size
}

func (fs *FileInfoMock) Sys() interface{} {
	return fs.sys
}
