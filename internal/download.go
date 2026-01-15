package internal

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
)

func FetchAsset(url, hashSpec, outPath, pkgName string) error {
	if strings.HasPrefix(url, "git://") {
		if fi, err := os.Stat(outPath); err == nil && fi.IsDir() {
			fmt.Printf("[*] %s git repo already cloned at %s, skipping clone\n", pkgName, outPath)
			return nil
		}

		repoURL := "https://" + strings.TrimPrefix(url, "git://")
		fmt.Printf("[*] cloning git repo %s into %s\n", repoURL, outPath)

		_, err := git.PlainClone(outPath, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
		if err != nil {
			return fmt.Errorf("git clone failed: %w", err)
		}
		return nil
	}

	if fi, err := os.Stat(outPath); err == nil && !fi.IsDir() {
		matches, err := VerifyFileHash(outPath, hashSpec)
		if err != nil {
			return fmt.Errorf("hash check failed for existing file %s: %w", outPath, err)
		}
		if matches {
			fmt.Printf("[*] %s file exists and hash matches, skipping download\n", pkgName)
			return nil
		}
		fmt.Printf("[*] %s file exists but hash mismatch, re-downloading\n", pkgName)
	}

	return downloadAndVerify(url, hashSpec, outPath)
}

func VerifyFileHash(path, hashSpec string) (bool, error) {
	algo, expected := parseHash(hashSpec)
	if algo == "" {
		return true, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	h, err := selectHash(algo)
	if err != nil {
		return false, err
	}

	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}

	actual := hex.EncodeToString(h.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(actual), []byte(strings.ToLower(expected))) == 1, nil
}

func downloadAndVerify(url, hashSpec, outPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error: %s", resp.Status)
	}

	tmp := outPath + ".part"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer out.Close()

	var h hash.Hash
	algo, expected := parseHash(hashSpec)
	if algo != "" {
		h, err = selectHash(algo)
		if err != nil {
			return err
		}
	}

	var writer io.Writer = out
	if h != nil {
		writer = io.MultiWriter(out, h)
	}

	if _, err := io.Copy(writer, resp.Body); err != nil {
		return err
	}

	if h != nil {
		actual := hex.EncodeToString(h.Sum(nil))
		if subtle.ConstantTimeCompare(
			[]byte(actual),
			[]byte(strings.ToLower(expected)),
		) != 1 {
			return fmt.Errorf("hash mismatch (%s): expected %s, got %s", algo, expected, actual)
		}
	}

	return os.Rename(tmp, outPath)
}

func parseHash(spec string) (algo, value string) {
	if spec == "" {
		return "", ""
	}

	parts := strings.SplitN(spec, ":", 2)
	if len(parts) != 2 || parts[1] == "none" {
		return "", ""
	}

	return strings.ToLower(parts[0]), parts[1]
}

func selectHash(algo string) (hash.Hash, error) {
	switch algo {
	case "md5":
		return md5.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("unsupported hash: %s", algo)
	}
}