package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
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

// DirEmpty returns whether given directory is empty or not.
// Folder is considered empty if only the given ignore files are present.
func DirEmpty(path string, ignore ...string) (bool, error) {
	f, err := os.Open(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return true, nil
	} else if err != nil {
		return false, err
	}

	defer f.Close()

	names, err := f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	if len(names) == 0 {
		return true, nil
	}

	for _, name := range names {
		var shouldIgnore bool
		for _, iName := range ignore {
			if name == iName {
				shouldIgnore = true
			}
		}

		if !shouldIgnore {
			return false, nil
		}
	}

	return true, nil
}

// Create file creates file and it's path if not exists.
func CreateFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0o770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func CopyFile(srcFile, dstFile string) error {
	input, err := os.ReadFile(srcFile)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", srcFile, err)
	}

	if err := os.WriteFile(dstFile, input, 0o644); err != nil {
		return fmt.Errorf("failed to write to file %q: %w", dstFile, err)
	}

	return nil
}

func Unzip(source, fileName, outDir string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	destination := filepath.Join(outDir, fileName)

	for _, f := range reader.File {
		if f.Name != fileName {
			continue
		}

		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	destinationFile, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}

func CpAndChmodxFile(source, destination string) error {
	if err := CopyFile(source, destination); err != nil {
		return fmt.Errorf("failed to copy file %q to %q: %w", source, destination, err)
	}

	if err := os.Chmod(destination, 0o700); err != nil {
		return fmt.Errorf("failed to chmod 0700 file %q: %w", destination, err)
	}

	return nil
}
