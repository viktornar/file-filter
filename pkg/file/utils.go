package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	appendMode = os.O_RDWR | os.O_CREATE | os.O_APPEND
	truncateMode  = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	permissions       = 0644
)

func CreateFolder(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return os.Mkdir(path, os.ModePerm)
	}

	return nil
}

func IsFile(path string) bool {
	f, err := os.Stat(path)
	
	if err != nil {
		return false
	}
	
	return !f.IsDir()
}

func CopyFile(in string, out string) error {
	fin, err := os.Open(in)
	
	if err != nil {
		return fmt.Errorf("can't open a file: %v", err)
	}
	
	defer fin.Close()

	err = CreateFolder(filepath.Dir(out))
	if err != nil {
		return fmt.Errorf("error while creating output directory: %v", err)
	}

	fout, err := os.Create(out)
	
	if err != nil {
		return fmt.Errorf("can't create a new file: %v", err)
	}
	
	defer fout.Close()

	if _, err = io.Copy(fout, fin); err != nil {
		return fmt.Errorf("can't copy file contents: %v", err)
	}

	return nil
}

func OverwriteFile(path string, data []byte) error {
	f, err := os.OpenFile(path, truncateMode, permissions)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func ReadFile(path string) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, permissions)
	
	if err != nil {
		return nil, err
	}
	
	defer f.Close()

	res, err := io.ReadAll(f)
	
	if err != nil {
		return nil, err
	}

	return res, nil
}

func OpenFile(path string) (*os.File, error) {
	return os.OpenFile(path, appendMode, permissions)
}

func DeleteFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return os.Remove(path)
}

func FileScanner(file *os.File, cb func(string)) {
	if file == nil || cb == nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cb(scanner.Text())
	}
}
