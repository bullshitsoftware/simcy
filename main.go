package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/bullshitsoftware/simcy/network"
	"github.com/bullshitsoftware/simcy/storage"
	"github.com/gen2brain/go-unarr"
)

type Network interface {
	GetLastReleaseName() (string, error)
	DownloadRelease(release, destDir string) (string, error)
}

type Storage interface {
	HasRelease(release string) bool
	ExtractRelease(filename string) error
	CleanupReleases(keepRelease string) error
	Run(release string) error
}

var (
	storageDir = flag.String(
		"storage",
		func() string {
			if appData, found := os.LookupEnv("LOCALAPPDATA"); found {
				return path.Join(appData, "Simcy")
			}
			return path.Join(os.TempDir(), "simcy")
		}(),
		"where to store simcraft",
	)
	downloadsUrl = flag.String(
		"downloads-url",
		"http://downloads.simulationcraft.org/nightly/",
		"simcraft downloads page",
	)
)

func main() {
	flag.Parse()

	if err := os.MkdirAll(*storageDir, 0700); err != nil {
		fmt.Printf("[ERROR] Failed to create storage directory %s: %v\n", *storageDir, err)
		os.Exit(1)
	}

	os.Exit(run(
		network.NewNetwork(http.DefaultClient, *downloadsUrl),
		storage.NewStorage(
			func(path string) (storage.Archive, error) { return unarr.NewArchive(path) },
			func(name string, arg ...string) storage.Starter { return exec.Command(name, arg...) },
			*storageDir,
		),
	))
}

func run(network Network, storage Storage) int {
	fmt.Printf("[INFO] Storage directory is %s\n", *storageDir)

	fmt.Printf("[INFO] Checking last release version\n")
	lastRelease, err := network.GetLastReleaseName()
	if err != nil {
		fmt.Printf("[ERROR] Failed to get last release info, %v\n", err)
		return 1
	}

	if !storage.HasRelease(lastRelease) {
		fmt.Printf("[INFO] Downloading release %s\n", lastRelease)
		archive, err := network.DownloadRelease(lastRelease, *storageDir)
		if err != nil {
			fmt.Printf("[ERROR] Failed to download, %v\n", err)
			return 1
		}

		fmt.Printf("[INFO] Extracting...\n")
		if err = storage.ExtractRelease(archive); err != nil {
			fmt.Printf("[ERROR] Failed to extract, %v\n", err)
			return 1
		}

		if err = storage.CleanupReleases(lastRelease); err != nil {
			fmt.Printf("[WARNING] Failed to cleanup, %v\n", err)
		}
	} else {
		fmt.Printf("[INFO] Has last release\n")
	}

	fmt.Printf("[INFO] Running...\n")
	if err = storage.Run(lastRelease); err != nil {
		fmt.Printf("[ERROR] Failed to run, %v\n", err)
		return 1
	}

	return 0
}
