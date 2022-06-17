package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/google/go-github/v43/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

const (
	repositoryOwner      = "vegaprotocol"
	vegaRepository       = "vega"
	vegaVersion          = "v0.51.1"
	vegaWalletRepository = "vegawallet"
	vegaWalletVersion    = "v0.15.1"
	dataNodeRepository   = "data-node"
	dataNodeVersion      = "v0.51.1"
	vegaBinName          = "vega"
	vegaWalletBinName    = "vegawallet"
	dataNodeBin          = "data-node"
)

var (
	githubToken string
	installPath string
)

var installBinariesCmd = &cobra.Command{
	Use:   "install-bins",
	Short: "Automatically download and install supported versions of vega, vegawallet and data-node binaries.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(githubToken) == 0 {
			githubToken = os.Getenv("GITHUB_TOKEN")

			if len(githubToken) == 0 {
				return fmt.Errorf("--github-token flag or GITHUB_TOKEN must be defined")
			}
		}

		installPath, err := getInstallPath(installPath)
		if err != nil {
			return err
		}

		info, err := os.Lstat(installPath)
		if err != nil {
			return fmt.Errorf("failed to get info about install-path %q: %w", installPath, err)
		}

		if !info.IsDir() {
			return fmt.Errorf("install-path should be a should be a directory")
		}

		if err := installDependencies(githubToken, installPath); err != nil {
			return fmt.Errorf("failed to install dependencies: %w", err)
		}

		if err := utils.BinariesAccessible(vegaBinName, vegaWalletBinName, dataNodeBin); err != nil {
			return fmt.Errorf("failed to lookup installed binaries, please check %q is in $PATH: %w", installPath, err)
		}

		return nil
	},
}

func getInstallPath(installPath string) (string, error) {
	if len(installPath) != 0 {
		return installPath, nil
	}

	installPath = os.Getenv("GOBIN")
	if len(installPath) == 0 {
		return "", fmt.Errorf("GOBIN environment variable has not been found - please set install-path flag instead")
	}

	return installPath, nil
}

func init() {
	installBinariesCmd.PersistentFlags().StringVar(&githubToken,
		"github-token",
		"",
		"Github personal token. Uses GITHUB_TOKEN environment variable by default.",
	)
	installBinariesCmd.PersistentFlags().StringVar(&installPath,
		"install-path",
		"",
		"Install path for the binaries. Uses GOBIN environment variable by default.",
	)
	installBinariesCmd.MarkFlagRequired("github-token")
}

func installDependencies(githubToken, installPath string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		vegaAssetName := fmt.Sprintf("%s-%s-amd64", vegaBinName, runtime.GOOS)

		if err := downloadReleaseAsset(ctx, client, repositoryOwner, vegaRepository, vegaVersion, vegaAssetName, homePath); err != nil {
			return fmt.Errorf("failed to download binary for vega: %w", err)
		}

		vegaAssetPath := path.Join(homePath, vegaAssetName)
		if err := cpAndChmodxFile(
			vegaAssetPath,
			path.Join(installPath, vegaBinName),
		); err != nil {
			return err
		}

		os.Remove(vegaAssetPath)

		return nil
	})

	eg.Go(func() error {
		vegaWalletAssetName := fmt.Sprintf("%s-%s-%s.zip", vegaWalletBinName, runtime.GOOS, runtime.GOARCH)
		vegaWalletAssetPath := path.Join(homePath, vegaWalletAssetName)

		if err := downloadReleaseAsset(ctx, client, repositoryOwner, vegaWalletRepository, vegaWalletVersion, vegaWalletAssetName, homePath); err != nil {
			return fmt.Errorf("failed to download binary for vega wallet: %w", err)
		}

		log.Printf("Unziping %q from %q", vegaWalletBinName, vegaWalletAssetPath)

		if err := utils.Unzip(vegaWalletAssetPath, vegaWalletBinName, homePath); err != nil {
			return fmt.Errorf("failed to unzip file %q from %q: %w", vegaWalletBinName, vegaWalletAssetName, err)
		}

		vegaWalletBinaryPath := path.Join(homePath, vegaWalletBinName)
		if err := cpAndChmodxFile(
			vegaWalletBinaryPath,
			path.Join(installPath, vegaWalletBinName),
		); err != nil {
			return err
		}

		os.Remove(vegaWalletAssetPath)
		os.Remove(vegaWalletBinaryPath)

		return nil
	})

	eg.Go(func() error {
		dataNodeAsset := fmt.Sprintf("%s-%s-amd64", dataNodeBin, runtime.GOOS)

		if err := downloadReleaseAsset(ctx, client, repositoryOwner, dataNodeRepository, dataNodeVersion, dataNodeAsset, homePath); err != nil {
			return fmt.Errorf("failed to download binary for data-node: %w", err)
		}

		dataNodeAssetPath := path.Join(homePath, dataNodeAsset)
		if err := cpAndChmodxFile(
			dataNodeAssetPath,
			path.Join(installPath, dataNodeBin),
		); err != nil {
			return err
		}

		os.Remove(dataNodeAssetPath)

		return nil
	})

	return eg.Wait()
}

func cpAndChmodxFile(source, destination string) error {
	if err := utils.CpAndChmodxFile(source, destination); err != nil {
		return err
	}
	log.Printf("Successfully copied from %q to %q", source, destination)
	return nil
}

func downloadReleaseAsset(ctx context.Context, client *github.Client, owner, repository, releaseTag, assetName, downloadDir string) error {
	log.Printf("Downloading release asset for %q with tag %q", repository, releaseTag)

	releases, resp, err := client.Repositories.ListReleases(ctx, owner, repository, nil)
	if err != nil {
		return err
	}

	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Github Token Expiration: %v\n", resp.TokenExpiration)
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

	assets, _, err := client.Repositories.ListReleaseAssets(ctx, owner, repository, releaseID, nil)
	if err != nil {
		return err
	}

	var assetID int64
	for _, asset := range assets {
		if asset.GetName() == assetName {
			assetID = asset.GetID()
		}
	}

	if assetID == 0 {
		return fmt.Errorf("asset %q in repository %q not found", repository, assetName)
	}

	ra, _, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repository, assetID, http.DefaultClient)
	if err != nil {
		return fmt.Errorf("failed to download release asset: %w", err)
	}
	defer ra.Close()

	downloadPath := path.Join(downloadDir, assetName)

	file, err := utils.CreateFile(downloadPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", downloadPath, err)
	}
	defer file.Close()

	all, err := ioutil.ReadAll(ra)
	if err != nil {
		return fmt.Errorf("failed to read  %q: %w", downloadPath, err)
	}

	if _, err = file.Write(all); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	log.Printf("Asset for %q with tag %q successfully downloaded to %q", repository, releaseTag, downloadPath)

	return nil
}
