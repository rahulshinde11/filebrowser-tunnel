package internal

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetCacheDir returns the cache directory for storing binaries
func GetCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, ".cache", "filebrowser-tunnel")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cacheDir, nil
}

// GetFreePort returns an available port by binding to :0
func GetFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("failed to find free port: %w", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// ExtractTunnelURL extracts the trycloudflare.com URL from cloudflared output
// Uses simple string operations instead of regexp to reduce binary size
func ExtractTunnelURL(output string) string {
	const suffix = ".trycloudflare.com"
	idx := strings.Index(output, suffix)
	if idx == -1 {
		return ""
	}

	// Find the start of the URL (look for https://)
	start := strings.LastIndex(output[:idx], "https://")
	if start == -1 {
		return ""
	}

	// Extract URL and find where it ends
	url := output[start:]
	end := strings.Index(url[len("https://"):], suffix) + len("https://") + len(suffix)
	if end > len(url) {
		end = len(url)
	}

	return url[:end]
}

// GetPlatformInfo returns the current OS and architecture
func GetPlatformInfo() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}

// getFilebrowserURL returns the download URL for filebrowser binary
func getFilebrowserURL(goos, goarch string) string {
	// Format: https://github.com/filebrowser/filebrowser/releases/latest/download/{os}-{arch}-filebrowser.tar.gz
	return fmt.Sprintf(
		"https://github.com/filebrowser/filebrowser/releases/latest/download/%s-%s-filebrowser.tar.gz",
		goos, goarch,
	)
}

// getCloudflaredURL returns the download URL for cloudflared binary
func getCloudflaredURL(goos, goarch string) string {
	// Format varies by OS:
	// Linux: https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-{arch}
	// Darwin: https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-darwin-{arch}.tgz
	if goos == "darwin" {
		return fmt.Sprintf(
			"https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-%s-%s.tgz",
			goos, goarch,
		)
	}
	return fmt.Sprintf(
		"https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-%s-%s",
		goos, goarch,
	)
}

// ClearCache removes all cached binaries
func ClearCache() error {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(cacheDir); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	fmt.Println("Cache cleared successfully")
	return nil
}
