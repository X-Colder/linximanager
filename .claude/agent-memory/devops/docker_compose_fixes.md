---
name: Docker Compose 部署已知坑点与修复
description: docker compose up 能正常运行所依赖的关键决策和非显而易见的修复点
type: project
---

以下问题在初始化 Docker Compose 部署时发现并已修复：

**1. PostgreSQL 镜像必须用 postgis/postgis:16-3.4-alpine**
迁移 SQL (001_init_schema.sql) 第一行执行 `CREATE EXTENSION IF NOT EXISTS postgis`，官方 postgres:16-alpine 不含 postgis，会导致初始化失败。
**Why:** 地理位置字段（商家坐标）依赖 postgis 扩展。
**How to apply:** 所有引用 postgres 镜像的地方都用 postgis/postgis:16-3.4-alpine，不要换回 postgres:16-alpine。

**2. REDIS_ADDR 变量格式为 "host:port"**
config.go 的 RedisConfig.Addr 字段读取的是 `REDIS_ADDR`（格式 "redis:6379"），不是分开的 REDIS_HOST + REDIS_PORT。旧的 deploy/env/.env.development 用了分开格式，与 config.go 不符。
**Why:** config.go 直接赋值给 redis.Options.Addr。
**How to apply:** 在 docker-compose.yml 的 backend 服务始终设置 `REDIS_ADDR=redis:6379`，不要用 REDIS_HOST+REDIS_PORT。

**3. admin.Dockerfile 的 nginx.conf 必须在 build context 内**
admin.Dockerfile 的 build context 是 `./frontend/admin`，所以 COPY 只能引用该目录内的文件。nginx 配置文件放在 `frontend/admin/nginx.conf`，而不是 `deploy/nginx/admin-nginx.conf`（后者在 build context 之外，无法 COPY）。
**How to apply:** admin 镜像的内嵌 nginx 配置永远是 `frontend/admin/nginx.conf`。

**4. ai-engine 目录是必须存在的**
docker-compose.yml 以 `./ai-engine` 为 build context，若目录不存在则 compose up 直接报错。已创建最小化 Python FastAPI 骨架（app/main.py、requirements.txt 等）。

**5. docker-compose.prod.yml 不能用 `ports: !reset []` 语法**
该 YAML 扩展语法在标准 docker compose 中不受支持，用 `ports: []` 覆盖即可。
