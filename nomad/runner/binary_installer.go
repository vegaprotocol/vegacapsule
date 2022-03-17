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

const (
	nomadBinaryVersion = "1.2.6"
	osDarwin           = "darwin"
	arm64Arch          = "arm64"
	amd64Arch          = "amd64"
)

var (
	nomadBinaryName = fmt.Sprintf("nomad_%s", nomadBinaryVersion)
)

func nomadBinaryPath() (string, error) {
	homeDir, err := capsuleHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".vegacapsule", nomadBinaryName), nil
}

func nomadDownloadUrl(binaryVersion, goos, arch string) (string, error) {
	url := fmt.Sprintf("https://releases.hashicorp.com/nomad/%s/nomad_%s_%s_%s.zip", binaryVersion, binaryVersion, goos, arch)

	resp, err := http.Head(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		return url, nil
	}

	// For Mac M1 fallback to amd64 if original architecture not found
	if goos == osDarwin && arch == arm64Arch {
		return nomadDownloadUrl(binaryVersion, goos, amd64Arch)
	}

	return "", fmt.Errorf("failed to get existing url for nomad binary %s: %w", binaryVersion, err)
}

func installNomadBinary(binPath string) error {
	url, err := nomadDownloadUrl(nomadBinaryVersion, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return fmt.Errorf("failed downloading nomad binary: expected version of nomad not found: %w", err)
	}

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
