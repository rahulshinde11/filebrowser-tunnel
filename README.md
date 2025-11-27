# filebrowser-tunnel

Expose any directory via a public URL with one command. Combines [filebrowser](https://github.com/filebrowser/filebrowser) with [Cloudflare Tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/) (trycloudflare) to create instant, secure file sharing.

## Quick Start

### One-liner (Hosted on GitHub Pages)

The easiest way to install:

```bash
curl -sL https://rahulshinde11.github.io/filebrowser-tunnel | sh
```

Or with a specific directory:

```bash
curl -sL https://rahulshinde11.github.io/filebrowser-tunnel | sh -s -- /path/to/share
```

### Direct Download

Download the binary for your platform from the [releases page](https://github.com/rahulshinde11/filebrowser-tunnel/releases/latest) and run:

```bash
./filebrowser-tunnel              # Serve current directory
./filebrowser-tunnel /path/to/dir # Serve specific directory
```

## Features

- **Zero configuration** - Just run and get a public URL
- **Automatic binary management** - Downloads filebrowser and cloudflared automatically
- **Cached binaries** - First run downloads, subsequent runs are instant
- **Cross-platform** - Linux (amd64, arm64) and macOS (Intel & Apple Silicon)
- **Secure tunnel** - Uses Cloudflare's free trycloudflare.com service
- **Web-based file browser** - Full filebrowser UI for browsing and downloading files

## Usage

```bash
filebrowser-tunnel [options] [directory]

Arguments:
  directory    Directory to serve (default: current directory)

Options:
  --help       Show help message
  --version    Show version
  --clean      Clear cached binaries and exit

Examples:
  filebrowser-tunnel                    # Serve current directory
  filebrowser-tunnel /path/to/dir       # Serve specific directory
  filebrowser-tunnel ~/Downloads        # Serve Downloads folder
  filebrowser-tunnel --clean            # Clear cached binaries
```

## How It Works

1. **Binary Detection** - Checks `~/.cache/filebrowser-tunnel/` for cached binaries
2. **Auto Download** - If not cached, downloads filebrowser and cloudflared for your platform
3. **Start Filebrowser** - Launches filebrowser on a random available port (no auth mode)
4. **Create Tunnel** - Starts cloudflared quick tunnel pointing to the local filebrowser
5. **Display URL** - Shows the public `*.trycloudflare.com` URL

## Project Structure

```
filebrowser-tunnel/
├── main.go                      # Entry point, CLI parsing
├── internal/
│   ├── downloader.go            # Binary download and caching
│   ├── runner.go                # Process management
│   └── utils.go                 # Utilities (port, URL, platform)
├── scripts/
│   ├── build.sh                 # Cross-platform build script
│   ├── publish.sh               # Docker build and push
│   ├── install.sh.template      # curl | sh install template
│   └── docker-entrypoint.sh     # Docker entrypoint
├── Dockerfile
├── Makefile
├── go.mod
├── .gitignore
├── .dockerignore
└── README.md
```

## Building from Source

### Prerequisites

- Go 1.22+
- Docker (for cross-compilation)
- UPX (optional, for binary compression)

### Build Commands

```bash
# Build for current platform
go build -o filebrowser-tunnel .

# Build for all platforms
make build-all

# Or use the build script (includes UPX compression)
./scripts/build.sh
```

### Build Targets

```bash
make help          # Show all available targets
make build-all     # Build for all platforms
make linux         # Build Linux binaries only
make darwin        # Build macOS binaries only
make clean         # Clean build artifacts
make run           # Build and run locally
```

## Self-Hosting

### GitHub Pages (Free Static Hosting)

This project is automatically deployed to GitHub Pages on every release!

**No setup required** - Just use the hosted version:

```bash
curl -sL https://rahulshinde11.github.io/filebrowser-tunnel/ | sh
```

**How it works:**
- Push a tag → GitHub Actions deploys to GitHub Pages
- Uses a "polyglot" file that works as both HTML (browsers) and shell script (curl)
- Browsers see a landing page, curl gets the install script

**To enable for your fork:**
1. Go to Settings → Pages
2. Source: "GitHub Actions"
3. Your distribution will be at `https://yourusername.github.io/filebrowser-tunnel/`

---

### Docker Container (GHCR)

Host your own distribution server using GitHub Container Registry:

#### Using Pre-built Image from GHCR

```bash
# Pull and run the latest version
docker run -d \
  -p 80:80 \
  -e DOMAIN=https://your-domain.com \
  ghcr.io/rahulshinde11/filebrowser-tunnel:latest

# Or use a specific version
docker run -d \
  -p 80:80 \
  -e DOMAIN=https://your-domain.com \
  ghcr.io/rahulshinde11/filebrowser-tunnel:v1.0.0
```

#### Building and Pushing to Your Own Registry

```bash
# Build for all platforms
make build-all

# Build Docker image
make docker-image

# Push to your own GHCR (or Docker Hub)
docker tag filebrowser-tunnel:latest ghcr.io/yourusername/filebrowser-tunnel:latest
docker push ghcr.io/yourusername/filebrowser-tunnel:latest
```

---

### Automatic Releases (CI/CD)

This project uses GitHub Actions for automatic releases:

- **Push to `main`** → Automatically deploys Docker image as `:main` tag to GHCR
- **Push a tag** (e.g., `v1.0.0`) → Creates GitHub Release with binaries + deploys versioned Docker images **+ deploys to GitHub Pages**

```bash
# Create and push release
git tag v1.0.0
git push origin v1.0.0
```

### User Access

Users can then run:

```bash
curl -sL https://your-domain.com | sh
```



## Cache Location

Binaries are cached in `~/.cache/filebrowser-tunnel/`:
- `filebrowser` - The filebrowser binary
- `cloudflared` - The cloudflared binary

Clear the cache with:

```bash
filebrowser-tunnel --clean
```

## Security Notes

- The filebrowser instance runs in **no-auth mode** for convenience
- The tunnel URL is randomly generated and not guessable
- For sensitive data, consider adding authentication or using a different solution
- Cloudflare tunnel traffic is encrypted end-to-end

## Supported Platforms

| OS | Architecture | Status |
|----|--------------|--------|
| Linux | amd64 | ✅ |
| Linux | arm64 | ✅ |
| macOS | amd64 (Intel) | ✅ |
| macOS | arm64 (Apple Silicon) | ✅ |

## License

MIT License

## Credits

- [filebrowser](https://github.com/filebrowser/filebrowser) - Web-based file manager
- [cloudflared](https://github.com/cloudflare/cloudflared) - Cloudflare Tunnel client
