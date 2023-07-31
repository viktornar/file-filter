package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupUtilsTests(t testing.TB) (string, func()) {
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

func TestFilePathWalkDir(t *testing.T) {
	testDir, teardown := setupUtilsTests(t)
	defer teardown()

	files, err := FilePathWalkDir(testDir)

	if err != nil {
		t.Fatal("expected to walk through given path")
	}

	expectedNames := []string{"file.txt", "file_1.txt", "file_2.txt", "file_3.txt"}

	for idx, file := range files {
		if !strings.Contains(file, expectedNames[idx]) {
			t.Errorf("Expected path %s to contain file name %s", file, expectedNames[idx])
		}
	}
}
