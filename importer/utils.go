package importer

import (
	"fmt"
	"io/ioutil"
)

func createTempFile(content string) (string, error) {
	file, err := ioutil.TempFile("", "vegacapsule_import")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	if _, err := file.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write content to temp file: %w", err)
	}

	return file.Name(), nil
}
