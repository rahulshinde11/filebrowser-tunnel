#!/bin/bash
set -e

# Generate index.html and install.sh for GitHub Pages

REPO="${1:-username/repo}"
VERSION="${2:-latest}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DOMAIN="https://${REPO%%/*}.github.io/${REPO##*/}"
DOWNLOAD_DOMAIN="https://github.com/${REPO}/releases/download/${VERSION}"

echo "Generating files for repo: ${REPO}, version: ${VERSION}"
echo "Pages domain: ${DOMAIN}"
echo "Download domain: ${DOWNLOAD_DOMAIN}"

# Generate index.html from template
REPO_URL="https://github.com/${REPO}"
if [ -f "${SCRIPT_DIR}/index.html.template" ]; then
    sed -e "s|{{DOMAIN}}|${DOMAIN}/install.sh|g" \
        -e "s|{{REPO_URL}}|${REPO_URL}|g" \
        "${SCRIPT_DIR}/index.html.template" > index.html
    echo "Generated index.html"
else
    echo "Error: index.html.template not found in ${SCRIPT_DIR}"
    exit 1
fi

# Generate install.sh from template
if [ -f "${SCRIPT_DIR}/install.sh.template" ]; then
    sed "s|{{DOMAIN}}|${DOWNLOAD_DOMAIN}|g" "${SCRIPT_DIR}/install.sh.template" > install.sh
    echo "Generated install.sh"
else
    echo "Error: install.sh.template not found in ${SCRIPT_DIR}"
    exit 1
fi

echo "Done! Files generated:"
echo "  - index.html (HTML page for browsers)"
echo "  - install.sh (Shell script for curl | sh)"
