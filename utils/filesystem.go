package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
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

	if err := os.Chmod(destination, 0700); err != nil {
		return fmt.Errorf("failed to chmod 0700 file %q: %w", destination, err)
	}

	return nil
}
