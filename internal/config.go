package internal

import (
	"errors"
	"fmt"
	"os"
	"net/url"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Packages map[string]Package `yaml:"packages"`
}

type Package struct {
	Version        string           `yaml:"version"`
	Tags           []string         `yaml:"tags,omitempty"`
	GlobalAssets   map[string]Asset `yaml:"global_assets,omitempty"`
	PlatformAssets map[string]Asset `yaml:"platform_assets,omitempty"`
}

type Asset struct {
	URL  string `yaml:"url"`
	Hash string `yaml:"hash"`
}

func validateURL(rawurl string) error {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	switch strings.ToLower(parsed.Scheme) {
	case "http", "https", "git":
		if parsed.Scheme == "git" {
			if parsed.Path == "" {
				return errors.New("git URL must have a repository path")
			}
		} else {
			if parsed.Host == "" {
				return errors.New("URL missing host")
			}
		}
		return nil
	default:
		return errors.New("unsupported URL scheme: " + parsed.Scheme)
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if len(cfg.Packages) == 0 {
		return nil, errors.New("no packages defined")
	}

	// Validate hashes on load
	for pkgName, pkg := range cfg.Packages {
		for assetKey, asset := range pkg.GlobalAssets {
			if err := validateHash(asset.Hash); err != nil {
				return nil, fmt.Errorf("package %s global asset %s: %w", pkgName, assetKey, err)
			}
			if err := validateURL(asset.URL); err != nil {
				return nil, fmt.Errorf("package %s global asset %s: invalid URL: %w", pkgName, assetKey, err)
			}
		}
		for assetKey, asset := range pkg.PlatformAssets {
			if err := validateHash(asset.Hash); err != nil {
				return nil, fmt.Errorf("package %s platform asset %s: %w", pkgName, assetKey, err)
			}
			if err := validateURL(asset.URL); err != nil {
				return nil, fmt.Errorf("package %s platform asset %s: invalid URL: %w", pkgName, assetKey, err)
			}
		}
	}

	return &cfg, nil
}

func validateHash(hashSpec string) error {
	if hashSpec == "" || hashSpec == "sha256:none" || hashSpec == "md5:none" || hashSpec == "sha512:none" {
		return nil // "none" means no hash checking
	}

	algo, value := parseHash(hashSpec)
	if algo == "" || value == "" {
		return fmt.Errorf("invalid hash format: %s", hashSpec)
	}

	if !isHex(value) {
		return fmt.Errorf("hash value is not valid hex: %s", value)
	}

	// Validate length matches algo
	switch algo {
	case "md5":
		if len(value) != 32 {
			return fmt.Errorf("md5 hash length invalid: %s", value)
		}
	case "sha256":
		if len(value) != 64 {
			return fmt.Errorf("sha256 hash length invalid: %s", value)
		}
	case "sha512":
		if len(value) != 128 {
			return fmt.Errorf("sha512 hash length invalid: %s", value)
		}
	default:
		return fmt.Errorf("unsupported hash algorithm: %s", algo)
	}

	return nil
}

func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') ||
			(c >= 'a' && c <= 'f') ||
			(c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}