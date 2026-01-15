package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/ctrlaltleet/loadout/internal"
)

func main() {
	var (
		configPath  string
		outputDir   string
		selectStr   string
		listOnly    bool
		concurrency int
		platform    string
	)

	flag.StringVar(&configPath, "config", "", "path to config file (required)")
	flag.StringVar(&outputDir, "output-dir", "./downloads", "output directory (default ./downloads)")
	flag.StringVar(&selectStr, "select", "", "comma-separated package names and/or tags to select (or 'all')")
	flag.BoolVar(&listOnly, "list", false, "list packages for current platform")
	flag.IntVar(&concurrency, "concurrency", 4, "number of parallel downloads (default 4)")
	flag.StringVar(&platform, "platform", "", "override detected platform (e.g. linux_amd64)")
	flag.Parse()

	if flag.NFlag() == 0 && flag.NArg() == 0 {
		fmt.Println(`loadout

@ctrlaltleet
`)		
		flag.Usage()
		os.Exit(0)
	}

	if configPath == "" {
		internal.Fatal("missing -config")
	}
	if concurrency < 1 {
		internal.Fatal("concurrency must be at least 1")
	}

	cfg, err := internal.LoadConfig(configPath)
	internal.FatalIf(err)

	if platform == "" {
		platform = runtime.GOOS + "_" + runtime.GOARCH
	}

	if listOnly {
		internal.ListPackages(cfg, platform, selectStr)
		return
	}

	if selectStr == "" {
		internal.Fatal("must specify -select for download or use -list")
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		internal.FatalIf(err)
	}

	gitDir := filepath.Join(outputDir, "git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		internal.FatalIf(err)
	}

	// Create platform subdirectory for platform-specific assets
	platformDir := filepath.Join(outputDir, platform)
	if err := os.MkdirAll(platformDir, 0o755); err != nil {
		internal.FatalIf(err)
	}

	jobs := make(chan internal.Job)
	results := make(chan internal.Result)

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go internal.Worker(&wg, jobs, results)
	}

	go func() {
		defer close(jobs)
		filters := internal.ParseSelect(selectStr)
		selectedPackages := internal.FilterPackages(cfg.Packages, filters)

		if len(selectedPackages) == 0 {
			fmt.Println("No packages matched the selection/filter criteria.")
			return
		}

		for name, pkg := range selectedPackages {
			for gaKey, asset := range pkg.GlobalAssets {
				var outPath string
				if strings.HasPrefix(asset.URL, "git://") {
					outPath = filepath.Join(gitDir, name)
				} else {
					outPath = filepath.Join(outputDir, filepath.Base(asset.URL))
				}
				jobs <- internal.Job{
					Name:    fmt.Sprintf("%s (global:%s)", name, gaKey),
					Version: pkg.Version,
					URL:     asset.URL,
					Hash:    asset.Hash,
					OutPath: outPath,
				}
			}

			if asset, ok := pkg.PlatformAssets[platform]; ok {
				outPath := filepath.Join(platformDir, filepath.Base(asset.URL))
				jobs <- internal.Job{
					Name:    fmt.Sprintf("%s (platform:%s)", name, platform),
					Version: pkg.Version,
					URL:     asset.URL,
					Hash:    asset.Hash,
					OutPath: outPath,
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var failed int
	for res := range results {
		if res.Err != nil {
			fmt.Printf("[!] %s failed: %v\n", res.Name, res.Err)
			failed++
		} else {
			fmt.Printf("[+] %s fetched successfully\n", res.Name)
		}
	}

	if failed > 0 {
		os.Exit(1)
	}
}
