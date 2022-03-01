package nomad

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	nomadBinaryVersion     = "1.2.5"
	nomadBinaryName        = fmt.Sprintf("nomad_%s", nomadBinaryVersion)
	defaultNomadBinaryPath = filepath.Join(build.Default.GOPATH, "bin", nomadBinaryName)
)

func nomadBinaryPath() string {
	// TODO change this to os.Executable path instead
	goBin := os.Getenv("GOBIN")
	if goBin == "" {
		return defaultNomadBinaryPath
	}

	return filepath.Join(goBin, nomadBinaryName)
}

func installNomadBinary() error {
	url := fmt.Sprintf("https://releases.hashicorp.com/nomad/%s/nomad_%s_%s_%s.zip", nomadBinaryVersion, nomadBinaryVersion, runtime.GOOS, runtime.GOARCH)

	log.Printf("Nomad binary was not found in. Installing from: %s", url)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zReader, err := zip.NewReader(bytes.NewReader(body), resp.ContentLength)
	if err != nil {
		return err
	}

	file, err := zReader.Open("nomad")
	if err != nil {
		return err
	}

	binPath := nomadBinaryPath()

	binFile, err := os.Create(binPath)
	if err != nil {
		return err
	}

	err = os.Chmod(binPath, 0755)
	if err != nil {
		return err
	}

	if _, err := io.Copy(binFile, file); err != nil {
		return err
	}

	log.Printf("Nomad binary was successfully installed in: %q", binPath)

	return nil
}

var nomadConfigTemplate = `
plugin "docker" {
	config {
		auth {
			%s
		}
	}
}
`

func templateConfig() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf(nomadConfigTemplate, `helper = "osxkeychain"`), nil
	case "windows":
		// TODO fix this for Windows machine
		return fmt.Sprintf(nomadConfigTemplate, `helper = "osxkeychain"`), nil
	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		authConfig := fmt.Sprintf(`config = "%s"`, filepath.Join(homeDir, ".docker", "config.json"))
		return fmt.Sprintf(nomadConfigTemplate, authConfig), nil
	default:
		return "", fmt.Errorf("platform not supported")
	}
}

func generateConfig() (string, error) {
	conf, err := templateConfig()
	if err != nil {
		return "", err
	}

	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(filepath.Dir(execPath), "nomad-config.hcl")

	if err := ioutil.WriteFile(configPath, []byte(conf), 0644); err != nil {
		return "", err
	}

	return configPath, nil
}

func StartAgent() error {
	nomadBinary := nomadBinaryPath()

	if _, err := exec.LookPath(nomadBinary); err != nil {
		if err := installNomadBinary(); err != nil {
			return fmt.Errorf("failed to install nomad binary: %w", err)
		}
	}

	configPath, err := generateConfig()
	if err != nil {
		return fmt.Errorf("failed to generate nomad config: %w", err)
	}

	args := []string{"agent", "-dev", "-bind", "0.0.0.0", "-config", configPath}

	log.Printf("Starting nomad with with: %s %v", nomadBinary, args)

	command := exec.Command(nomadBinary, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("failed to execute binary %s with error: %s", nomadBinary, err.Error())
	}

	return nil
}
