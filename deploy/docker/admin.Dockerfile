# ===== Stage 1: Builder =====
FROM node:20-alpine3.19 AS builder

WORKDIR /app

COPY package.json package-lock.json* pnpm-lock.yaml* yarn.lock* ./
RUN if [ -f pnpm-lock.yaml ]; then \
      npm install -g pnpm && pnpm install --frozen-lockfile; \
    elif [ -f yarn.lock ]; then \
      yarn install --frozen-lockfile; \
    else \
      npm ci; \
    fi

COPY . .
RUN npm run build

# ===== Stage 2: Runner =====
FROM nginx:1.25-alpine

RUN addgroup -S nginxgroup && \
    adduser -S nginxuser -G nginxgroup && \
    mkdir -p /var/cache/nginx /var/run && \
    chown -R nginxuser:nginxgroup /var/cache/nginx /var/log/nginx /var/run

COPY --from=builder /app/dist /usr/share/nginx/html
COPY deploy/nginx/admin-nginx.conf /etc/nginx/conf.d/default.conf

RUN chown -R nginxuser:nginxgroup /usr/share/nginx/html

USER nginxuser

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:80/ || exit 1
