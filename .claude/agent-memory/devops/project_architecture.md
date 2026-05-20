---
name: 灵犀掌柜项目架构概览
description: 灵犀掌柜项目的技术栈、服务拓扑和部署架构，用于部署方案决策
type: project
---

灵犀掌柜（linximanager）是一个智能库存决策平台，采用模块化单体架构。

**技术栈**
- 后端：Go 1.22 + Gin，入口 `backend/cmd/server/main.go`，监听 8080
- AI 引擎：Python FastAPI + gRPC，端口 8000（HTTP）/ 50051（gRPC）
- 管理后台：Vue3 + TypeScript + Element Plus + Vite，位于 `frontend/admin/`
- 小程序：uni-app，商家端 `frontend/merchant-app/`，C 端 `frontend/consumer-app/`
- 数据库：PostgreSQL 16（JSONB、分区表），Redis 7，MinIO，Meilisearch v1.6

**服务端口映射（Docker Compose 开发环境）**
- backend: 8080, admin: 3000, ai-engine: 8000/50051
- postgres: 5432, redis: 6379, minio: 9000/9001, meilisearch: 7700

**生产拓扑**：K8s 集群，backend HPA 3-10 副本（CPU>70%触发），ai-engine 2副本固定
- Ingress：admin.linximanager.com → admin-service，linximanager.com/api/ → backend-service

**Why:** 团队需要从开发到生产的完整部署方案，支持灰度发布和自动扩缩容。
**How to apply:** 提出部署相关问题时，优先参考已有的 deploy/ 目录配置；镜像仓库地址为 registry.cn-hangzhou.aliyuncs.com/linximanager/。
