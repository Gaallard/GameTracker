#!/bin/sh
set -e

# Use default value if BACKEND_URL is not set
export BACKEND_URL="${BACKEND_URL:-http://localhost:8080}"

echo "Processing nginx template with BACKEND_URL=${BACKEND_URL}"

# Process nginx template with environment variables
if [ -f /etc/nginx/templates/default.conf.template ]; then
    envsubst '${BACKEND_URL}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf
    # Remove template to prevent double processing by original entrypoint
    rm /etc/nginx/templates/default.conf.template
    echo "Nginx configuration file created successfully"
else
    echo "Warning: Template file not found at /etc/nginx/templates/default.conf.template"
fi

# Execute the original nginx entrypoint (which handles other initialization)
exec /docker-entrypoint.sh.original "$@"

