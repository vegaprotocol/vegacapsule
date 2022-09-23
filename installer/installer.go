package installer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/blang/semver"
	"github.com/google/go-github/v43/github"
	"golang.org/x/sync/errgroup"
)

const (
	repository      = "vega"
	repositoryOwner = "vegaprotocol"

	vegaBinName     = "vega"
	walletBinName   = "vegawallet"
	dataNodeBinName = "data-node"
)

var (
	minSupportedVersion = semver.MustParse("0.54.0")
	assetsToInstall     = map[string]string{
		formatAssetName(vegaBinName):     vegaBinName,
		formatAssetName(walletBinName):   walletBinName,
		formatAssetName(dataNodeBinName): dataNodeBinName,
	}
)

type InstalledBins map[string]string

func (ib InstalledBins) lookup(name string) (string, bool) {
	s, ok := ib[name]
	return s, ok
}

func (ib InstalledBins) VegaPath() (string, bool) {
	return ib.lookup(vegaBinName)
}

func (ib InstalledBins) WalletPath() (string, bool) {
	return ib.lookup(walletBinName)
}

func (ib InstalledBins) DataNodePath() (string, bool) {
	return ib.lookup(dataNodeBinName)
}

type Installer struct {
	repository      string
	repositoryOwner string
	binDirectory    string
	installPath     string
	client          *github.Client
}

func New(binDirectory, installPath string) *Installer {
	return &Installer{
		repository:      repository,
		repositoryOwner: repositoryOwner,
		binDirectory:    binDirectory,
		installPath:     installPath,
		client:          github.NewClient(nil),
	}
}

type asset struct {
	ID         int64
	Name       string
	BinaryName string
}

func (i Installer) getAssets(ctx context.Context, releaseTag string, requestedAssets map[string]string) ([]asset, error) {
	log.Printf("Downloading release asset for %q with tag %q", i.repository, releaseTag)

	releases, resp, err := i.client.Repositories.ListReleases(ctx, i.repositoryOwner, i.repository, nil)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("release in repository %q with tag %q not found", i.repository, releaseTag)
	}

	ghAssets, _, err := i.client.Repositories.ListReleaseAssets(ctx, i.repositoryOwner, i.repository, releaseID, nil)
	if err != nil {
		return nil, err
	}

	assets := []asset{}
	for _, a := range ghAssets {
		if binaryName, ok := requestedAssets[a.GetName()]; ok {
			assets = append(assets, asset{
				ID:         a.GetID(),
				Name:       a.GetName(),
				BinaryName: binaryName,
			})
		}
	}

	if len(assets) == 0 {
		return nil, fmt.Errorf("node assets %v in repository %q not found", requestedAssets, i.repository)
	}

	return assets, nil
}

func (i Installer) downloadAsset(ctx context.Context, asset asset, releaseTag string) (string, error) {
	ra, _, err := i.client.Repositories.DownloadReleaseAsset(ctx, i.repositoryOwner, i.repository, asset.ID, http.DefaultClient)
	if err != nil {
		return "", fmt.Errorf("failed to download release asset: %w", err)
	}
	defer ra.Close()

	downloadPath := path.Join(i.binDirectory, asset.Name)
	binPath := path.Join(i.binDirectory, asset.BinaryName)

	file, err := utils.CreateFile(downloadPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %q: %w", downloadPath, err)
	}
	defer file.Close()

	all, err := ioutil.ReadAll(ra)
	if err != nil {
		return "", fmt.Errorf("failed to read  %q: %w", downloadPath, err)
	}

	if _, err = file.Write(all); err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	log.Printf("Asset for %q with tag %q successfully downloaded to %q", i.repository, releaseTag, downloadPath)

	log.Printf("Unziping %q from %q to %q", asset.BinaryName, asset.Name, binPath)

	if err := utils.Unzip(downloadPath, asset.BinaryName, i.binDirectory); err != nil {
		return "", fmt.Errorf("failed to unzip file %q from %q: %w", asset.Name, asset.BinaryName, err)
	}

	// Make sure the file has executable perms
	if err := os.Chmod(binPath, 0700); err != nil {
		return "", fmt.Errorf("failed to chmod 0700 file %q: %w", binPath, err)
	}

	log.Printf("Successfully unzipped %q from %q to %q", asset.BinaryName, asset.Name, binPath)

	defer os.Remove(downloadPath)

	if i.installPath != "" {
		if err := cpAndChmodxFile(binPath, i.installPath); err != nil {
			return "", fmt.Errorf("failed to copy binary to predefined path %q: %w", i.installPath, err)
		}
	}

	return binPath, nil
}

func (i Installer) Install(ctx context.Context, releaseTag string) (InstalledBins, error) {
	// Parse in semver without the "v" prefix
	releaseVersion := strings.TrimLeft(releaseTag, "v")
	v, err := semver.Parse(releaseVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version from relase tag %v", releaseTag)
	}

	if v.LT(minSupportedVersion) {
		return nil, fmt.Errorf("requested version %q must be bigger or equal then minimum supported version %q", v, minSupportedVersion)
	}

	log.Printf("Starting to install binaries to %q", i.binDirectory)

	downloadAssets, err := i.getAssets(ctx, releaseTag, assetsToInstall)
	if err != nil {
		return nil, err
	}

	var mut sync.Mutex
	installedBinsPaths := InstalledBins{}

	eg, ctx := errgroup.WithContext(ctx)
	for _, asset := range downloadAssets {
		asset := asset

		eg.Go(func() error {
			binPath, err := i.downloadAsset(ctx, asset, releaseTag)
			if err != nil {
				return fmt.Errorf("failed to download %s: %w", asset.Name, err)
			}

			mut.Lock()
			installedBinsPaths[asset.BinaryName] = binPath
			mut.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	log.Printf("Successfully to install binaries to %q", i.binDirectory)

	return installedBinsPaths, nil
}
