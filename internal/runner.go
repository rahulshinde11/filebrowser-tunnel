package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

// ProcessManager manages the filebrowser and cloudflared processes
type ProcessManager struct {
	filebrowserCmd *exec.Cmd
	cloudflaredCmd *exec.Cmd
	tunnelURL      string
	mu             sync.Mutex
	urlChan        chan string
}

// NewProcessManager creates a new ProcessManager
func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		urlChan: make(chan string, 1),
	}
}

// StartFilebrowser starts the filebrowser process
func (pm *ProcessManager) StartFilebrowser(binaryPath string, port int, directory string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Get absolute path for the directory
	absDir, err := getAbsolutePath(directory)
	if err != nil {
		return fmt.Errorf("failed to resolve directory: %w", err)
	}

	// Verify directory exists
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", absDir)
	}

	pm.filebrowserCmd = exec.Command(
		binaryPath,
		"--noauth",
		"--address", "0.0.0.0",
		"--port", fmt.Sprintf("%d", port),
		"--root", absDir,
	)

	// Capture output for debugging
	pm.filebrowserCmd.Stdout = os.Stdout
	pm.filebrowserCmd.Stderr = os.Stderr

	if err := pm.filebrowserCmd.Start(); err != nil {
		return fmt.Errorf("failed to start filebrowser: %w", err)
	}

	fmt.Printf("🗂️  Filebrowser started on port %d (serving: %s)\n", port, absDir)
	return nil
}

// StartCloudflared starts the cloudflared tunnel
func (pm *ProcessManager) StartCloudflared(binaryPath string, localPort int) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.cloudflaredCmd = exec.Command(
		binaryPath,
		"tunnel",
		"--url", fmt.Sprintf("http://localhost:%d", localPort),
	)

	// Create pipe for stderr to capture the tunnel URL
	stderr, err := pm.cloudflaredCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := pm.cloudflaredCmd.Start(); err != nil {
		return fmt.Errorf("failed to start cloudflared: %w", err)
	}

	// Read stderr in background to find the tunnel URL
	go pm.parseCloudflaredOutput(stderr)

	return nil
}

// parseCloudflaredOutput reads cloudflared stderr and extracts the tunnel URL
func (pm *ProcessManager) parseCloudflaredOutput(stderr io.ReadCloser) {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()

		// Look for the tunnel URL
		if url := ExtractTunnelURL(line); url != "" {
			pm.mu.Lock()
			pm.tunnelURL = url
			pm.mu.Unlock()

			// Send URL to channel (non-blocking)
			select {
			case pm.urlChan <- url:
			default:
			}
		}
	}
}

// WaitForTunnelURL waits for the tunnel URL to be available
func (pm *ProcessManager) WaitForTunnelURL(timeout time.Duration) (string, error) {
	select {
	case url := <-pm.urlChan:
		return url, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timeout waiting for tunnel URL")
	}
}

// GetTunnelURL returns the current tunnel URL
func (pm *ProcessManager) GetTunnelURL() string {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.tunnelURL
}

// Stop stops both processes gracefully
func (pm *ProcessManager) Stop() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	fmt.Println("\n🛑 Shutting down...")

	if pm.cloudflaredCmd != nil && pm.cloudflaredCmd.Process != nil {
		pm.cloudflaredCmd.Process.Kill()
		pm.cloudflaredCmd.Wait()
	}

	if pm.filebrowserCmd != nil && pm.filebrowserCmd.Process != nil {
		pm.filebrowserCmd.Process.Kill()
		pm.filebrowserCmd.Wait()
	}

	fmt.Println("✓ Stopped")
}

// Wait waits for both processes to finish
func (pm *ProcessManager) Wait() error {
	var wg sync.WaitGroup
	var filebrowserErr, cloudflaredErr error

	if pm.filebrowserCmd != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filebrowserErr = pm.filebrowserCmd.Wait()
		}()
	}

	if pm.cloudflaredCmd != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cloudflaredErr = pm.cloudflaredCmd.Wait()
		}()
	}

	wg.Wait()

	if filebrowserErr != nil {
		return fmt.Errorf("filebrowser error: %w", filebrowserErr)
	}
	if cloudflaredErr != nil {
		return fmt.Errorf("cloudflared error: %w", cloudflaredErr)
	}

	return nil
}

// getAbsolutePath returns the absolute path of a directory
func getAbsolutePath(path string) (string, error) {
	if path == "" || path == "." {
		return os.Getwd()
	}

	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = home + path[1:]
	}

	return absPath(path)
}

// absPath returns the absolute path
func absPath(path string) (string, error) {
	if path[0] == '/' {
		return path, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return cwd + "/" + path, nil
}

