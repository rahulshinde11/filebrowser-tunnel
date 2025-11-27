#!/bin/sh
set -e

if [ -z "$DOMAIN" ]; then
    echo "Error: DOMAIN environment variable is not set."
    echo "Usage: docker run -e DOMAIN=https://your-domain.com -p 80:80 yourname/filebrowser-tunnel"
    exit 1
fi

# Remove trailing slash if present to avoid double slashes
DOMAIN=$(echo "$DOMAIN" | sed 's:/*$::')

echo "Generating pages with domain: $DOMAIN"

# Generate install.sh (served to curl/wget)
sed "s|{{DOMAIN}}|${DOMAIN}|g" /usr/share/nginx/html/install.sh.template > /usr/share/nginx/html/install.sh

# Generate index.html (served to browsers)
sed "s|{{DOMAIN}}|${DOMAIN}|g" /usr/share/nginx/html/index.html.template > /usr/share/nginx/html/index.html

# Execute the CMD (nginx)
exec "$@"
