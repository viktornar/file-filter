package file

import (
	"file-filter/pkg/slice"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupWatcherTests(t testing.TB) (string, func()) {
	testDir, err := os.MkdirTemp(".", "")
	if err != nil {
		t.Fatal(err)
	}

	err = WriteFile(filepath.Join(testDir, "file.txt"), []byte{})
	if err != nil {
		t.Fatal(err)
	}

	files := []string{"file_1.txt", "file_2.txt", "file_3.txt"}

	for _, f := range files {
		filePath := filepath.Join(testDir, f)
		if err := WriteFile(filePath, []byte{}); err != nil {
			t.Fatal(err)
		}
	}

	testDirTwo := filepath.Join(testDir, "testDirTwo")
	CreateDir(testDirTwo)

	err = WriteFile(filepath.Join(testDirTwo, "file_recursive.txt"),
		[]byte{})
	if err != nil {
		t.Fatal(err)
	}

	abs, err := filepath.Abs(testDir)
	if err != nil {
		os.RemoveAll(testDir)
		t.Fatal(err)
	}

	return abs, func() {
		if os.RemoveAll(testDir); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileInfo(t *testing.T) {
	modTime := time.Now()

	fileInfo := &fileInfo{
		name:    "fileInfo",
		size:    1,
		mode:    os.ModeDir,
		modTime: modTime,
		sys:     nil,
		dir:     true,
	}

	if fileInfo.Name() != "fileInfo" {
		t.Fatalf("expected fileInfo.Name() to be 'fileInfo', got %s", fileInfo.Name())
	}
	if fileInfo.IsDir() != true {
		t.Fatalf("expected fileInfo.IsDir() to be true, got %t", fileInfo.IsDir())
	}
	if fileInfo.Size() != 1 {
		t.Fatalf("expected fileInfo.Size() to be 1, got %d", fileInfo.Size())
	}
	if fileInfo.Sys() != nil {
		t.Fatalf("expected fileInfo.Sys() to be nil, got %v", fileInfo.Sys())
	}
	if fileInfo.ModTime() != modTime {
		t.Fatalf("expected fileInfo.ModTime() to be %v, got %v", modTime, fileInfo.ModTime())
	}
	if fileInfo.Mode() != os.ModeDir {
		t.Fatalf("expected fileInfo.Mode() to be os.ModeDir, got %#v", fileInfo.Mode())
	}

	w := NewWatcher()

	w.wg.Done()

	go func() {
		w.TriggerEvent(Create, fileInfo)
	}()

	e := <-w.Event

	if e.FileInfo != fileInfo {
		t.Fatal("expected e.FileInfo to be equal to fileInfo")
	}
}

func TestWatcherAdd(t *testing.T) {
	testDir, teardown := setupWatcherTests(t)
	defer teardown()

	w := NewWatcher()

	err := w.Add("-")
	if err == nil {
		t.Error("expected error to not be nil")
	}

	if err := w.Add(testDir); err != nil {
		t.Fatal(err)
	}

	if len(w.files) != 7 {
		t.Errorf("expected len(w.files) to be 7, got %d", len(w.files))
	}

	if slice.IndexOf[string](w.names, testDir) == -1 {
		t.Errorf("expected w.names to contain testDir")
	}

	if _, found := w.files[testDir]; !found {
		t.Errorf("expected to find %s", testDir)
	}

	if w.files[testDir].Name() != filepath.Base(testDir) {
		t.Errorf("expected w.files[%q].Name() to be %s, got %s",
			testDir, testDir, w.files[testDir].Name())
	}

	fileRecursive := filepath.Join(testDir, "testDirTwo", "file_recursive.txt")
	if _, found := w.files[fileRecursive]; !found {
		t.Errorf("expected to find %s", fileRecursive)
	}

	fileTxt := filepath.Join(testDir, "file.txt")
	if _, found := w.files[fileTxt]; !found {
		t.Errorf("expected to find %s", fileTxt)
	}

	if w.files[fileTxt].Name() != "file.txt" {
		t.Errorf("expected w.files[%q].Name() to be file.txt, got %s",
			fileTxt, w.files[fileTxt].Name())
	}

	dirTwo := filepath.Join(testDir, "testDirTwo")
	if _, found := w.files[dirTwo]; !found {
		t.Errorf("expected to find %s directory", dirTwo)
	}

	if w.files[dirTwo].Name() != "testDirTwo" {
		t.Errorf("expected w.files[%q].Name() to be testDirTwo, got %s",
			dirTwo, w.files[dirTwo].Name())
	}
}
