# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Accept build arguments for environment variables
# If VITE_API_URL is not provided, use empty string for relative paths (nginx will proxy)
ARG VITE_API_URL=""
ENV VITE_API_URL=$VITE_API_URL

# Build the application
RUN npm run build

# Production stage
FROM nginx:alpine

# Install gettext for envsubst
RUN apk add --no-cache gettext

# Copy built files to nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration template
COPY nginx.conf /etc/nginx/templates/default.conf.template

# Save original entrypoint and copy custom entrypoint script
RUN cp /docker-entrypoint.sh /docker-entrypoint.sh.original
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Set default backend URL if not provided
ENV BACKEND_URL=http://localhost:8080

# Expose port
EXPOSE 80

# Use custom entrypoint that processes template before starting nginx
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
