package bundletool

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// githubRelease represents the JSON structure of a GitHub release.
type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// EnsureBundleTool checks if bundletool exists, and downloads the latest version if it doesn't.
// It returns the absolute path to the jar.
func EnsureBundleTool() (string, error) {
	// 1. Check if user explicitly set an environment variable
	if path := os.Getenv("BUNDLETOOL_PATH"); path != "" {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// 2. Define the global user directory (~/.bundletool-mcp/bundletool.jar)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	bundleDir := filepath.Join(homeDir, ".bundletool-mcp")
	jarPath := filepath.Join(bundleDir, "bundletool.jar")

	// 3. Check if we already downloaded it previously
	if _, err := os.Stat(jarPath); err == nil {
		return jarPath, nil // Already exists!
	}

	// 4. Not found? Let's fetch the latest release info from GitHub.
	fmt.Fprintf(os.Stderr, "Bundletool not found locally. Fetching latest release info from GitHub...\n")
	downloadURL, err := getLatestReleaseURL()
	if err != nil {
		return "", fmt.Errorf("failed to get latest release URL: %v", err)
	}

	// 5. Download the jar
	fmt.Fprintf(os.Stderr, "Downloading bundletool from %s to %s...\n", downloadURL, jarPath)

	if err := os.MkdirAll(bundleDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to download bundletool: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status downloading bundletool: %s", resp.Status)
	}

	out, err := os.Create(jarPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed to write to file: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Download complete!\n")
	return jarPath, nil
}

func getLatestReleaseURL() (string, error) {
	apiURL := "https://api.github.com/repos/google/bundletool/releases/latest"
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".jar") {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no .jar asset found in latest release")
}
