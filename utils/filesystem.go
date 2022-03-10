package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CapsuleHome() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".vegacapsule"), nil
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// Create file creates file and it's path if not exists.
func CreateFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func CopyFile(srcFile, dstFile string) error {
	input, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", srcFile, err)
	}

	if err := ioutil.WriteFile(dstFile, input, 0644); err != nil {
		return fmt.Errorf("failed to write to file %q: %w", dstFile, err)
	}

	return nil
}
