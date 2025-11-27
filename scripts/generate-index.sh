#!/bin/bash
set -e

# Generate polyglot index file for GitHub Pages
# This file works as both HTML (for browsers) and shell script (for curl)

REPO="${1:-username/repo}"
VERSION="${2:-latest}"
BASE_URL="https://github.com/${REPO}/releases/download/${VERSION}"

cat > index.html << 'EOF'
#!/bin/sh
true <<'ENDOFHTML'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>filebrowser-tunnel - Instant File Sharing</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
            color: #333;
        }
        .container {
            background: white;
            border-radius: 16px;
            padding: 48px;
            max-width: 700px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
        }
        h1 {
            font-size: 2.5em;
            margin-bottom: 16px;
            color: #667eea;
        }
        .tagline {
            font-size: 1.2em;
            color: #666;
            margin-bottom: 32px;
        }
        .code-block {
            background: #1e1e1e;
            color: #d4d4d4;
            padding: 20px;
            border-radius: 8px;
            font-family: 'Monaco', 'Courier New', monospace;
            font-size: 0.9em;
            margin: 24px 0;
            overflow-x: auto;
            position: relative;
        }
        .code-block code {
            display: block;
        }
        .copy-btn {
            position: absolute;
            top: 12px;
            right: 12px;
            background: #667eea;
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 0.85em;
            transition: all 0.2s;
        }
        .copy-btn:hover {
            background: #5568d3;
            transform: translateY(-1px);
        }
        .copy-btn.copied {
            background: #48bb78;
        }
        .features {
            list-style: none;
            margin: 32px 0;
        }
        .features li {
            padding: 12px 0;
            padding-left: 32px;
            position: relative;
        }
        .features li:before {
            content: "✓";
            position: absolute;
            left: 0;
            color: #48bb78;
            font-weight: bold;
            font-size: 1.2em;
        }
        .links {
            margin-top: 32px;
            padding-top: 32px;
            border-top: 1px solid #e2e8f0;
        }
        .links a {
            color: #667eea;
            text-decoration: none;
            font-weight: 500;
            margin-right: 24px;
        }
        .links a:hover {
            text-decoration: underline;
        }
        .note {
            background: #fef3c7;
            border-left: 4px solid #f59e0b;
            padding: 16px;
            margin: 24px 0;
            border-radius: 4px;
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 filebrowser-tunnel</h1>
        <p class="tagline">Expose any directory via a public URL with one command</p>
        
        <h2 style="margin-top: 32px; margin-bottom: 16px; color: #333;">Quick Start</h2>
        <div class="code-block">
            <button class="copy-btn" onclick="copyCode(this)">Copy</button>
            <code id="install-cmd">curl -sL REPLACE_BASE_URL | sh</code>
        </div>

        <div class="note">
            <strong>Note:</strong> Or download binaries directly from the <a href="REPLACE_REPO_URL/releases/latest" target="_blank" style="color: #f59e0b; font-weight: bold;">releases page</a>
        </div>

        <h2 style="margin-top: 32px; margin-bottom: 16px; color: #333;">Features</h2>
        <ul class="features">
            <li>Zero configuration - Just run and get a public URL</li>
            <li>Automatic binary management - Downloads dependencies automatically</li>
            <li>Cross-platform - Linux & macOS (amd64, arm64)</li>
            <li>Secure tunnel - Uses Cloudflare's free service</li>
            <li>Web-based file browser - Full UI for browsing and downloading</li>
        </ul>

        <div class="links">
            <a href="REPLACE_REPO_URL" target="_blank">GitHub Repository</a>
            <a href="REPLACE_REPO_URL/releases" target="_blank">Releases</a>
            <a href="REPLACE_REPO_URL#readme" target="_blank">Documentation</a>
        </div>
    </div>

    <script>
        // Replace placeholders with actual values
        const REPO = 'REPLACE_REPO';
        const BASE_URL = window.location.href.replace(/\/$/, '');
        const REPO_URL = 'https://github.com/' + REPO;
        
        document.body.innerHTML = document.body.innerHTML
            .replace(/REPLACE_REPO_URL/g, REPO_URL)
            .replace(/REPLACE_BASE_URL/g, BASE_URL)
            .replace(/REPLACE_REPO/g, REPO);

        function copyCode(btn) {
            const code = document.getElementById('install-cmd').textContent;
            navigator.clipboard.writeText(code).then(() => {
                btn.textContent = 'Copied!';
                btn.classList.add('copied');
                setTimeout(() => {
                    btn.textContent = 'Copy';
                    btn.classList.remove('copied');
                }, 2000);
            });
        }
    </script>
</body>
</html>
ENDOFHTML

DOMAIN="REPLACE_DOMAIN"
BINARY_NAME="filebrowser-tunnel"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Map architectures
case "${ARCH}" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
esac

# Check for supported platforms
if [ "$OS" != "linux" ] && [ "$OS" != "darwin" ]; then
    echo "Error: Unsupported OS: $OS"
    exit 1
fi

if [ "$ARCH" != "amd64" ] && [ "$ARCH" != "arm64" ]; then
    echo "Error: Unsupported Architecture: $ARCH"
    exit 1
fi

FULL_BINARY_NAME="${BINARY_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL="${DOMAIN}/${FULL_BINARY_NAME}"

# Create a temp file
if [ -d "/tmp" ]; then
    TMP_FILE="/tmp/${BINARY_NAME}-$$"
else
    TMP_FILE="./${BINARY_NAME}-$$"
fi

# Cleanup function
cleanup() {
    rm -f "${TMP_FILE}"
}
trap cleanup EXIT

# Download
echo "Downloading ${FULL_BINARY_NAME} from ${DOMAIN}..."
if command -v curl >/dev/null 2>&1; then
    STATUS=$(curl -w "%{http_code}" -fsSL -o "${TMP_FILE}" "${DOWNLOAD_URL}")
    if [ "$STATUS" != "200" ] && [ ! -f "${TMP_FILE}" ]; then
        echo "Error: Failed to download binary (HTTP $STATUS)"
        exit 1
    fi
elif command -v wget >/dev/null 2>&1; then
    wget -qO "${TMP_FILE}" "${DOWNLOAD_URL}"
else
    echo "Error: curl or wget is required"
    exit 1
fi

chmod +x "${TMP_FILE}"

echo "Running ${BINARY_NAME}..."
echo "=========================="
"${TMP_FILE}" "$@"
EOF

# Replace placeholders
sed -i.bak "s|REPLACE_REPO|${REPO}|g" index.html
sed -i.bak "s|REPLACE_DOMAIN|https://github.com/${REPO}/releases/download/${VERSION}|g" index.html
rm -f index.html.bak

echo "Generated index.html"
