FROM node:20-alpine AS builder

WORKDIR /app

# Copy root workspace files
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/platform/package.json apps/platform/

# Install dependencies
RUN corepack enable && pnpm install --frozen-lockfile

# Copy platform source
COPY apps/platform/ apps/platform/

# Build
WORKDIR /app/apps/platform
RUN pnpm build

# ─── Runtime (Nginx) ─────────────────────
FROM nginx:alpine

COPY --from=builder /app/apps/platform/dist /usr/share/nginx/html

# SPA fallback
RUN echo 'server { \
    listen 80; \
    root /usr/share/nginx/html; \
    index index.html; \
    location / { \
        try_files $uri $uri/ /index.html; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
