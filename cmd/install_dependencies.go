package cmd

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"syscall"

	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/google/go-github/v43/github"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

const (
	owner                = "vegaprotocol"
	vegaRepository       = "vega"
	vegaWalletRepository = "vegawallet"
	dataNodeRepository   = "data-node"
)

var installCependenciesCmd = &cobra.Command{
	Use:   "install-deps",
	Short: "Automatically download and install supported versions of vega, vegawallet and data-node binaries.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// @TODO allow to set custom install PATH if gobin is missing
		goBin := os.Getenv("GOBIN")
		if len(goBin) == 0 {
			panic("missing GOBIN env variable")
		}

		info, err := os.Lstat(goBin)
		if err != nil {
			panic(err)
		}

		if !info.IsDir() {
			panic("GOBIN should be a directory")
		}

		fmt.Print("GitHub Token: ")
		byteToken, _ := terminal.ReadPassword(int(syscall.Stdin))
		println()
		token := string(byteToken)

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		eg, ctx := errgroup.WithContext(ctx)

		eg.Go(func() error {
			vegaBinaryName := fmt.Sprintf("vega-%s-amd64", runtime.GOOS)
			vegaBinaryPath := path.Join(homePath, vegaBinaryName)
			if err := downloadReleaseAsset(ctx, client, owner, vegaRepository, "v0.50.1", vegaBinaryName, homePath); err != nil {
				return fmt.Errorf("failed to download binary for vega: %w", err)
			}

			destination := path.Join(goBin, "vega")
			if err := utils.CopyFile(vegaBinaryPath, destination); err != nil {
				return fmt.Errorf("failed to copy file %q to %q: %w", vegaBinaryPath, destination, err)
			}

			if err := os.Chmod(destination, 0700); err != nil {
				return fmt.Errorf("failed to chmod 0700 file %q: %w", destination, err)
			}

			os.Remove(vegaBinaryName)

			return nil
		})

		eg.Go(func() error {
			vegaWalletAssetName := fmt.Sprintf("vegawallet-%s-%s.zip", runtime.GOOS, runtime.GOARCH)
			vegaWalletAssetPath := path.Join(homePath, vegaWalletAssetName)
			vegaWalletBinaryName := "vegawallet"
			vegaWalletBinaryPath := path.Join(homePath, vegaWalletBinaryName)

			if err := downloadReleaseAsset(ctx, client, owner, vegaWalletRepository, "v0.13.2", vegaWalletAssetName, homePath); err != nil {
				return fmt.Errorf("failed to download binary for vega wallet: %w", err)
			}

			if err := unzip(vegaWalletAssetPath, vegaWalletBinaryName, homePath); err != nil {
				return fmt.Errorf("failed to unzip file %q from %q: %w", vegaWalletBinaryName, vegaWalletAssetName, err)
			}

			destination := path.Join(goBin, vegaWalletBinaryName)
			if err := utils.CopyFile(vegaWalletBinaryPath, destination); err != nil {
				return fmt.Errorf("failed to copy file %q to %q: %w", vegaWalletBinaryPath, destination, err)
			}

			if err := os.Chmod(destination, 0700); err != nil {
				return fmt.Errorf("failed to chmod 0700 file %q: %w", vegaWalletBinaryPath, err)
			}

			os.Remove(vegaWalletAssetName)
			os.Remove(vegaWalletBinaryPath)

			return nil
		})

		eg.Go(func() error {
			dataNodeBinary := fmt.Sprintf("data-node-%s-amd64", runtime.GOOS)
			dataNodeBinaryPath := path.Join(homePath, dataNodeBinary)
			if err := downloadReleaseAsset(ctx, client, owner, dataNodeRepository, "v0.50.1", dataNodeBinary, homePath); err != nil {
				return fmt.Errorf("failed to download binary for data-node: %w", err)
			}

			destination := path.Join(goBin, "data-node")
			if err := utils.CopyFile(dataNodeBinaryPath, destination); err != nil {
				return fmt.Errorf("failed to copy file %q to %q: %w", dataNodeBinary, destination, err)
			}

			if err := os.Chmod(destination, 0700); err != nil {
				return fmt.Errorf("failed to chmod 0700 file %q: %w", destination, err)
			}

			os.Remove(dataNodeBinaryPath)

			return nil
		})

		return eg.Wait()
	},
}

func init() {
	// installCependenciesCmd.PersistentFlags().StringVar(&nomadConfigPath,
	// 	"nomad-config-path",
	// 	"",
	// 	"Allows to use Nomad configuration",
	// )
}

func downloadReleaseAsset(ctx context.Context, client *github.Client, owner, repository, releaseTag, assetName, downloadDir string) error {
	log.Printf("downloading release asset for %q with tag %q", repository, releaseTag)

	releases, _, err := client.Repositories.ListReleases(ctx, owner, repository, nil)
	if err != nil {
		return err
	}

	var releaseID int64
	for _, r := range releases {
		if r.GetTagName() == releaseTag {
			releaseID = r.GetID()
		}
	}

	if releaseID == 0 {
		return fmt.Errorf("release in repository %q with tag %q not found", repository, releaseTag)
	}

	assets, resp, err := client.Repositories.ListReleaseAssets(ctx, owner, repository, releaseID, nil)
	if err != nil {
		return err
	}

	var assetID int64
	for _, asset := range assets {
		if asset.GetName() == assetName {
			assetID = asset.GetID()
		}
	}

	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
	}

	if assetID == 0 {
		return fmt.Errorf("asset %q not found", assetName)
	}

	rc, _, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repository, assetID, http.DefaultClient)
	if err != nil {
		return fmt.Errorf("failed to download release asset: %w", err)
	}
	defer rc.Close()

	downloadPath := path.Join(downloadDir, assetName)

	file, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", downloadPath, err)
	}
	defer file.Close()

	all, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("failed to read  %q: %w", downloadPath, err)
	}

	_, err = file.Write(all)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	log.Printf("asset for %q with tag %q succesfully downloaded to %q", repository, releaseTag, downloadPath)

	return nil
}

func unzip(source, fileName, outDir string) error {
	log.Printf("unziping %q from %q", fileName, source)

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
