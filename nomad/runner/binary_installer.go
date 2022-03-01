package runner

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"code.vegaprotocol.io/vegacapsule/utils"
)

var (
	nomadBinaryVersion = "1.2.5"
	nomadBinaryName    = fmt.Sprintf("nomad_%s", nomadBinaryVersion)
)

func nomadBinaryPath() (string, error) {
	homeDir, err := capsuleHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".vegacapsule", nomadBinaryName), nil
}

func installNomadBinary(binPath string) error {
	url := fmt.Sprintf("https://releases.hashicorp.com/nomad/%s/nomad_%s_%s_%s.zip", nomadBinaryVersion, nomadBinaryVersion, runtime.GOOS, runtime.GOARCH)

	log.Printf("Nomad binary was not found in. Installing from: %s", url)

	c := http.Client{
		Timeout: time.Second * 20,
	}
	resp, err := c.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get nomad binary release from %q: %w", err, err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get nomad binary release with bad status: %q", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	zReader, err := zip.NewReader(bytes.NewReader(body), resp.ContentLength)
	if err != nil {
		return fmt.Errorf("failed to unzip nomad package: %w", err)
	}

	file, err := zReader.Open("nomad")
	if err != nil {
		return fmt.Errorf("failed to get nomad binary from unzipped folder: %w", err)
	}

	binFile, err := utils.CreateFile(binPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", binPath, err)
	}
	defer binFile.Close()

	err = os.Chmod(binPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to change permission for file %q: %w", binPath, err)
	}

	if _, err := io.Copy(binFile, file); err != nil {
		return fmt.Errorf("failed to copy content to file %q: %w", binPath, err)
	}

	log.Printf("Nomad binary was successfully installed at: %q", binPath)

	return nil
}
