package internal

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadProgress tracks download progress
type DownloadProgress struct {
	Total      int64
	Downloaded int64
	Writer     io.Writer
}

func (dp *DownloadProgress) Write(p []byte) (int, error) {
	n, err := dp.Writer.Write(p)
	dp.Downloaded += int64(n)
	dp.printProgress()
	return n, err
}

func (dp *DownloadProgress) printProgress() {
	if dp.Total > 0 {
		percent := float64(dp.Downloaded) / float64(dp.Total) * 100
		fmt.Printf("\r  Downloading... %.1f%% (%.2f MB / %.2f MB)",
			percent,
			float64(dp.Downloaded)/(1024*1024),
			float64(dp.Total)/(1024*1024))
	} else {
		fmt.Printf("\r  Downloading... %.2f MB", float64(dp.Downloaded)/(1024*1024))
	}
}

// downloadFile downloads a file from URL to the destination path
func downloadFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	progress := &DownloadProgress{
		Total:  resp.ContentLength,
		Writer: out,
	}

	_, err = io.Copy(progress, resp.Body)
	fmt.Println() // New line after progress
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// extractTarGz extracts a .tar.gz file and returns the path to the extracted binary
func extractTarGz(archivePath, destDir, binaryName string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	var extractedPath string

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read tar: %w", err)
		}

		// Look for the binary file (filebrowser or cloudflared)
		baseName := filepath.Base(header.Name)
		if header.Typeflag == tar.TypeReg && (baseName == binaryName || baseName == "filebrowser" || baseName == "cloudflared") {
			destPath := filepath.Join(destDir, binaryName)
			outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				return "", fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return "", fmt.Errorf("failed to extract file: %w", err)
			}
			outFile.Close()
			extractedPath = destPath
		}
	}

	if extractedPath == "" {
		return "", fmt.Errorf("binary not found in archive")
	}

	return extractedPath, nil
}

// ensureFilebrowser downloads and caches filebrowser if not present
func ensureFilebrowser(cacheDir string) (string, error) {
	goos, goarch := GetPlatformInfo()
	binaryPath := filepath.Join(cacheDir, "filebrowser")

	// Check if binary already exists
	if _, err := os.Stat(binaryPath); err == nil {
		return binaryPath, nil
	}

	fmt.Println("📦 Downloading filebrowser...")
	url := getFilebrowserURL(goos, goarch)

	// Download the archive
	archivePath := filepath.Join(cacheDir, "filebrowser.tar.gz")
	if err := downloadFile(url, archivePath); err != nil {
		return "", fmt.Errorf("failed to download filebrowser: %w", err)
	}

	// Extract the binary
	fmt.Println("  Extracting...")
	extractedPath, err := extractTarGz(archivePath, cacheDir, "filebrowser")
	if err != nil {
		return "", fmt.Errorf("failed to extract filebrowser: %w", err)
	}

	// Clean up archive
	os.Remove(archivePath)

	fmt.Println("  ✓ filebrowser ready")
	return extractedPath, nil
}

// ensureCloudflared downloads and caches cloudflared if not present
func ensureCloudflared(cacheDir string) (string, error) {
	goos, goarch := GetPlatformInfo()
	binaryPath := filepath.Join(cacheDir, "cloudflared")

	// Check if binary already exists
	if _, err := os.Stat(binaryPath); err == nil {
		return binaryPath, nil
	}

	fmt.Println("📦 Downloading cloudflared...")
	url := getCloudflaredURL(goos, goarch)

	if goos == "darwin" {
		// macOS: download .tgz archive
		archivePath := filepath.Join(cacheDir, "cloudflared.tgz")
		if err := downloadFile(url, archivePath); err != nil {
			return "", fmt.Errorf("failed to download cloudflared: %w", err)
		}

		// Extract the binary
		fmt.Println("  Extracting...")
		extractedPath, err := extractTarGz(archivePath, cacheDir, "cloudflared")
		if err != nil {
			return "", fmt.Errorf("failed to extract cloudflared: %w", err)
		}

		// Clean up archive
		os.Remove(archivePath)

		fmt.Println("  ✓ cloudflared ready")
		return extractedPath, nil
	}

	// Linux: direct binary download
	if err := downloadFile(url, binaryPath); err != nil {
		return "", fmt.Errorf("failed to download cloudflared: %w", err)
	}

	// Make executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return "", fmt.Errorf("failed to make cloudflared executable: %w", err)
	}

	fmt.Println("  ✓ cloudflared ready")
	return binaryPath, nil
}

// EnsureBinaries ensures both filebrowser and cloudflared are available
func EnsureBinaries() (filebrowserPath, cloudflaredPath string, err error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", "", err
	}

	// Check platform support
	goos, goarch := GetPlatformInfo()
	if goos != "linux" && goos != "darwin" {
		return "", "", fmt.Errorf("unsupported OS: %s", goos)
	}
	if goarch != "amd64" && goarch != "arm64" {
		return "", "", fmt.Errorf("unsupported architecture: %s", goarch)
	}

	fmt.Printf("🔍 Platform: %s/%s\n", goos, goarch)
	fmt.Printf("📁 Cache directory: %s\n\n", cacheDir)

	filebrowserPath, err = ensureFilebrowser(cacheDir)
	if err != nil {
		return "", "", err
	}

	cloudflaredPath, err = ensureCloudflared(cacheDir)
	if err != nil {
		return "", "", err
	}

	fmt.Println()
	return filebrowserPath, cloudflaredPath, nil
}

