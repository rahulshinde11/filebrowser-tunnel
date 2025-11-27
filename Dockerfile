FROM nginx:alpine

# Copy binaries from dist folder to web root
COPY dist/ /usr/share/nginx/html/

# Copy install script template
COPY scripts/install.sh.template /usr/share/nginx/html/install.sh.template

# Copy HTML template for browser visitors
COPY scripts/index.html.template /usr/share/nginx/html/index.html.template

# Copy custom nginx config for User-Agent detection
COPY scripts/nginx.conf /etc/nginx/conf.d/default.conf

# Copy entrypoint script
COPY scripts/docker-entrypoint.sh /docker-entrypoint.sh

# Make entrypoint executable
RUN chmod +x /docker-entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"]

# Default command to start Nginx
CMD ["nginx", "-g", "daemon off;"]
