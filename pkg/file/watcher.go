package file

import (
	"file-filter/pkg/logger"
	"file-filter/pkg/slice"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Op uint32

const (
	Create Op = iota
	Write
	Remove
	Rename
	Chmod
	Move
)

const (
	CreateName  = "CREATE"
	WriteName   = "WRITE"
	RemoveName  = "REMOVE"
	RenameName  = "RENAME"
	MoveName    = "MOVE"
	UnknownName = "UNKNOWN"
)

const (
	FilePathType = "FILE"
	DirPathType  = "DIRECTORY"
)

var Operations = map[Op]string{
	Create: CreateName,
	Write:  WriteName,
	Remove: RemoveName,
	Rename: RenameName,
	Move:   MoveName,
}

func (e Op) String() string {
	if op, found := Operations[e]; found {
		return op
	}
	return UnknownName
}

type Event struct {
	Op
	Path    string
	OldPath string
	os.FileInfo
}

func (e Event) String() string {
	if e.FileInfo == nil {
		return UnknownName
	}

	pathType := FilePathType
	if e.IsDir() {
		pathType = DirPathType
	}

	return fmt.Sprintf("%s %q %s [%s]", pathType, e.Name(), e.Op, e.Path)
}

type Watcher struct {
	Event   chan Event
	Error   chan error
	Closed  chan struct{}
	close   chan struct{}
	wg      *sync.WaitGroup
	mu      *sync.Mutex
	running bool
	names   []string
	files   map[string]os.FileInfo
}

func NewWatcher() *Watcher {
	var wg sync.WaitGroup
	wg.Add(1)

	return &Watcher{
		Event:  make(chan Event),
		Error:  make(chan error),
		Closed: make(chan struct{}),
		close:  make(chan struct{}),
		mu:     new(sync.Mutex),
		wg:     &wg,
		files:  make(map[string]os.FileInfo),
		names:  []string{},
	}
}

func (w *Watcher) Add(name string) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	name, err = filepath.Abs(name)
	if err != nil {
		return err
	}

	fileList, err := w.list(name)

	if err != nil {
		return err
	}

	for k, v := range fileList {
		w.files[k] = v
	}

	w.names = append(w.names, name)

	return nil
}

func (w *Watcher) list(name string) (map[string]os.FileInfo, error) {
	fileList := make(map[string]os.FileInfo)

	return fileList, filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error.Print(err)
			return err
		}

		fileList[path] = info
		logger.Debug.Printf("Adding file info to the file list %v\n", info)
		return nil
	})
}

func (w *Watcher) remove(name string) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	logger.Debug.Printf("Tryinh to remove file from list %s\n", name)

	name, err = filepath.Abs(name)
	if err != nil {
		logger.Error.Print(err)
		return err
	}

	slice.Remove[string](w.names, slice.IndexOf[string](w.names, name))

	info, found := w.files[name]

	if !found {
		return nil
	}

	if !info.IsDir() {
		delete(w.files, name)
		return nil
	}

	for path := range w.files {
		if strings.HasPrefix(path, name) {
			delete(w.files, path)
		}
	}
	return nil
}

func (w *Watcher) WatchedFiles() map[string]os.FileInfo {
	w.mu.Lock()
	defer w.mu.Unlock()

	files := make(map[string]os.FileInfo)
	for k, v := range w.files {
		files[k] = v
	}

	return files
}

func (w *Watcher) TriggerEvent(eventType Op, file os.FileInfo) {
	w.Wait()
	if file == nil {
		file = &FileInfoMock{name: "triggered event", modTime: time.Now()}
	}
	w.Event <- Event{Op: eventType, Path: "-", FileInfo: file}
}

func (w *Watcher) retrieveFileList() map[string]os.FileInfo {
	w.mu.Lock()
	defer w.mu.Unlock()

	fileList := make(map[string]os.FileInfo)

	var list map[string]os.FileInfo
	var err error

	for _, name := range w.names {
		list, err = w.list(name)
		if err != nil {
			if os.IsNotExist(err) {
				w.mu.Unlock()
				if name == err.(*os.PathError).Path {
					w.Error <- ErrWatchedFileDeleted
					w.remove(name)
				}
				w.mu.Lock()
			} else {
				w.Error <- err
			}
		}

		for k, v := range list {
			fileList[k] = v
		}
	}

	logger.Debug.Printf("Retrieved file list %v\n", fileList)

	return fileList
}

func (w *Watcher) Start(d time.Duration) error {
	logger.Debug.Printf("Starting watcher with duration %s\n", d)
	if d < time.Nanosecond {
		return ErrDurationTooShort
	}

	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return ErrWatcherRunning
	}
	w.running = true
	w.mu.Unlock()
	w.wg.Done()

	for {
		done := make(chan struct{})
		cancel := make(chan struct{})
		evt := make(chan Event)

		fileList := w.retrieveFileList()

		go func() {
			w.pollEvents(fileList, evt, cancel)
			done <- struct{}{}
		}()

	inner:
		for {
			select {
			case <-w.close:
				close(cancel)
				close(w.Closed)
				return nil
			case event := <-evt:
				w.Event <- event
			case <-done:
				break inner
			}
		}

		w.mu.Lock()
		w.files = fileList
		w.mu.Unlock()

		time.Sleep(d)
	}
}

func (w *Watcher) pollEvents(files map[string]os.FileInfo, evt chan Event, cancel chan struct{}) {
	w.mu.Lock()
	defer w.mu.Unlock()

	creates := make(map[string]os.FileInfo)
	removes := make(map[string]os.FileInfo)

	for path, info := range w.files {
		if _, found := files[path]; !found {
			removes[path] = info
		}
	}

	for path, info := range files {
		oldInfo, found := w.files[path]
		if !found {
			creates[path] = info
			continue
		}

		if oldInfo.ModTime() != info.ModTime() {
			select {
			case <-cancel:
				return
			case evt <- Event{Write, path, path, info}:
				logger.Debug.Printf("Sending evt %s", Event{Write, path, path, info})
			}
		}
	}

	for removePath, removeInfo := range removes {
		for createPath, createInfo := range creates {
			if SameFile(removeInfo, createInfo) {
				e := Event{
					Op:       Move,
					Path:     createPath,
					OldPath:  removePath,
					FileInfo: removeInfo,
				}

				if filepath.Dir(removePath) == filepath.Dir(createPath) {
					e.Op = Rename
				}

				delete(removes, removePath)
				delete(creates, createPath)

				select {
				case <-cancel:
					return
				case evt <- e:
					logger.Debug.Printf("Sending event %s", e)
				}
			}
		}
	}

	for path, info := range creates {
		select {
		case <-cancel:
			return
		case evt <- Event{Create, path, "", info}:
		}
	}

	for path, info := range removes {
		select {
		case <-cancel:
			return
		case evt <- Event{Remove, path, path, info}:
		}
	}
}

func (w *Watcher) Wait() {
	w.wg.Wait()
}

func (w *Watcher) Close() {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return
	}
	w.running = false
	w.files = make(map[string]os.FileInfo)
	w.names = []string{}
	w.mu.Unlock()
	w.close <- struct{}{}
}
