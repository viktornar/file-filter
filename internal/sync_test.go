package internal

import (
	"errors"
	"file-filter/pkg/file"
	"os"
	"path/filepath"
	"testing"
)

func setupSyncTests(t testing.TB) (string, func()) {
	testDir, err := os.MkdirTemp(".", "")
	if err != nil {
		t.Fatal(err)
	}

	err = file.WriteFile(filepath.Join(testDir, "file.txt"), []byte{})
	if err != nil {
		t.Fatal(err)
	}

	err = file.WriteFile(filepath.Join(testDir, "delete_file.txt"), []byte{})
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

func TestFileCreate(t *testing.T) {
	testDir, teardown := setupSyncTests(t)
	defer teardown()

	fileInfo := file.NewFileInfoMock("file.txt")

	event := file.Event{Op: file.Create, Path: "-", FileInfo: fileInfo}

	handleFileChange(event, testDir)

	if _, err := os.Stat(filepath.Join(testDir, event.Name())); errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected to find %s in %s directory", event.Name(), testDir)
	}
}

func TestFileRename(t *testing.T) {
	testDir, teardown := setupSyncTests(t)
	defer teardown()

	fileInfo := file.NewFileInfoMock("file.txt")
	newFileName := "new_file_name.txt"

	event := file.Event{Op: file.Rename, Path: testDir + "/" + newFileName, FileInfo: fileInfo}

	handleFileChange(event, testDir)

	if _, err := os.Stat(filepath.Join(testDir, newFileName)); errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected to find %s in %s directory", newFileName, testDir)
	}
}

func TestFileDelete(t *testing.T) {
	testDir, teardown := setupSyncTests(t)
	defer teardown()

	fileInfo := file.NewFileInfoMock("delete_file.txt")

	event := file.Event{Op: file.Remove, Path: "-", FileInfo: fileInfo}

	handleFileChange(event, testDir)

	if _, err := os.Stat(filepath.Join(testDir, event.Name())); err == nil {
		t.Errorf("expected to remove %s in %s directory", event.Name(), testDir)
	}
}
