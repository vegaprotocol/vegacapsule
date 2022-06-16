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
	"path"
	"path/filepath"
	"runtime"
	"time"

	"code.vegaprotocol.io/vegacapsule/utils"
)

const (
	nomadBinName       = "nomad"
	nomadBinaryVersion = "1.3.1"
)

var (
	nomadBinaryName = fmt.Sprintf("%s_%s", nomadBinName, nomadBinaryVersion)
)

func nomadBinaryPath() (string, error) {
	homeDir, err := capsuleHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, nomadBinaryName), nil
}

func nomadDownloadUrl(binaryVersion, goos, arch string) (string, error) {
	url := fmt.Sprintf("https://releases.hashicorp.com/nomad/%s/nomad_%s_%s_%s.zip", binaryVersion, binaryVersion, goos, arch)

	resp, err := http.Head(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		return url, nil
	}

	return "", fmt.Errorf("failed to get existing url for nomad binary %s: %w", binaryVersion, err)
}

func installNomadBinary(localBinPath, binInstallPath string) (err error) {
	url, err := nomadDownloadUrl(nomadBinaryVersion, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return fmt.Errorf("failed downloading nomad binary: expected version of nomad not found: %w", err)
	}

	log.Printf("Nomad binary was not found in. Installing from: %s", url)

	c := http.Client{
		Timeout: time.Second * 120,
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

	file, err := zReader.Open(nomadBinName)
	if err != nil {
		return fmt.Errorf("failed to get nomad binary from unzipped folder: %w", err)
	}
	defer file.Close()

	binFile, err := utils.CreateFile(localBinPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", localBinPath, err)
	}
	defer binFile.Close()

	err = os.Chmod(localBinPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to change permission for file %q: %w", localBinPath, err)
	}

	if _, err := io.Copy(binFile, file); err != nil {
		return fmt.Errorf("failed to copy content to file %q: %w", localBinPath, err)
	}

	defer func() {
		if err != nil {
			_ = os.Remove(binFile.Name())
		}
	}()

	log.Printf("Nomad binary was successfully installed at: %q", localBinPath)

	binInstallPathAbs := path.Join(binInstallPath, nomadBinName)

	if err := utils.CpAndChmodxFile(localBinPath, binInstallPathAbs); err != nil {
		return fmt.Errorf("failed to install Nomad to %q: %w", binInstallPath, err)
	}

	if err := utils.BinariesAccessible(nomadBinName); err != nil {
		return fmt.Errorf("failed to lookup installed binaries, please check %q is in $PATH: %w", binInstallPath, err)
	}

	log.Printf("Nomad binary was successfully installed at provided path: %q", binInstallPathAbs)

	return nil
}
