# =============================================================================
# 管理后台（Vue3 + Nginx）多阶段构建
#
# Build context: ./frontend/admin
# 命令示例:
#   docker build -f deploy/docker/admin.Dockerfile \
#                --build-arg VITE_API_BASE_URL=http://localhost:8080 \
#                -t linximanager/admin:latest ./frontend/admin
# =============================================================================

# ===== Stage 1: 前端构建 =====
FROM m.daocloud.io/docker.io/library/node:20-alpine3.19 AS builder

WORKDIR /app

# 接收构建时 API 地址（默认指向开发 backend，Nginx 反代则填 /api）
ARG VITE_API_BASE_URL=http://localhost:8080
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

# 先复制 lockfile 利用 Docker 层缓存；node_modules 只在依赖变化时重建
COPY package.json package-lock.json* pnpm-lock.yaml* yarn.lock* ./

RUN if [ -f pnpm-lock.yaml ]; then \
      npm install -g pnpm && pnpm install --frozen-lockfile; \
    elif [ -f yarn.lock ]; then \
      yarn install --frozen-lockfile; \
    else \
      npm ci; \
    fi

COPY . .

# TypeScript 类型检查 + Vite 生产构建
RUN npm run build

# ===== Stage 2: Nginx 静态文件服务 =====
FROM m.daocloud.io/docker.io/library/nginx:1.25-alpine

# 创建非 root 用户，nginx worker 进程无需 root
RUN addgroup -S nginxgroup && \
    adduser -S nginxuser -G nginxgroup && \
    mkdir -p /var/cache/nginx /var/run && \
    chown -R nginxuser:nginxgroup /var/cache/nginx /var/log/nginx /var/run

# 注意：Nginx 主进程需要以 root 绑定端口 80，
# 但 worker 进程会通过 nginx.conf 的 user 指令降权。
# 因此此处不切换到 nginxuser，由 nginx 自身完成降权。

# 从构建阶段复制 dist 产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 内嵌 SPA 的 Nginx 配置（此文件与 admin 镜像一同打包，不依赖挂载）
# 该文件位于 frontend/admin/ 目录内，与 build context 一致
COPY nginx.conf /etc/nginx/conf.d/default.conf

RUN chown -R nginxuser:nginxgroup /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:80/ping || exit 1
