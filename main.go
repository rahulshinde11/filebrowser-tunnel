package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shinde11/filebrowser-tunnel/internal"
)

var version = "dev"

func main() {
	// CLI flags
	showVersion := flag.Bool("version", false, "Show version")
	showHelp := flag.Bool("help", false, "Show help")
	cleanCache := flag.Bool("clean", false, "Clear cached binaries")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "filebrowser-tunnel - Expose directory via Cloudflare tunnel")
		fmt.Fprintln(os.Stderr, "\nUsage: filebrowser-tunnel [options] [directory]")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		fmt.Fprintln(os.Stderr, "  --help      Show help")
		fmt.Fprintln(os.Stderr, "  --version   Show version")
		fmt.Fprintln(os.Stderr, "  --clean     Clear cached binaries")
	}

	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("filebrowser-tunnel version %s\n", version)
		os.Exit(0)
	}

	// Handle help flag
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Handle clean flag
	if *cleanCache {
		if err := internal.ClearCache(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Get directory to serve
	directory := "."
	if flag.NArg() > 0 {
		directory = flag.Arg(0)
	}

	// Run the tunnel
	if err := run(directory); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(directory string) error {
	fmt.Println("=== filebrowser-tunnel ===")
	fmt.Println("Expose files via Cloudflare tunnel")

	// Ensure binaries are available
	filebrowserPath, cloudflaredPath, err := internal.EnsureBinaries()
	if err != nil {
		return fmt.Errorf("failed to ensure binaries: %w", err)
	}

	// Get a free port
	port, err := internal.GetFreePort()
	if err != nil {
		return fmt.Errorf("failed to get free port: %w", err)
	}

	// Create process manager
	pm := internal.NewProcessManager()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		pm.Stop()
		os.Exit(0)
	}()

	// Start filebrowser
	if err := pm.StartFilebrowser(filebrowserPath, port, directory); err != nil {
		return fmt.Errorf("failed to start filebrowser: %w", err)
	}

	// Wait a moment for filebrowser to start
	time.Sleep(1 * time.Second)

	// Start cloudflared tunnel
	fmt.Println("🌐 Starting Cloudflare tunnel...")
	if err := pm.StartCloudflared(cloudflaredPath, port); err != nil {
		pm.Stop()
		return fmt.Errorf("failed to start cloudflared: %w", err)
	}

	// Wait for tunnel URL
	url, err := pm.WaitForTunnelURL(30 * time.Second)
	if err != nil {
		pm.Stop()
		return fmt.Errorf("failed to get tunnel URL: %w", err)
	}

	// Display the URL
	fmt.Printf("\n🔗 Access: %s\n", url)
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Wait for processes to finish
	return pm.Wait()
}
