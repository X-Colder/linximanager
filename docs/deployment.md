# 灵犀掌柜 - 部署与运维指南

> 项目负责人：大哈 | 邮箱：915788160@qq.com | 更新日期：2026-05-20

---

## 目录

1. [前置条件](#1-前置条件)
2. [开发环境快速启动](#2-开发环境快速启动)
3. [生产环境部署](#3-生产环境部署)
4. [数据库迁移](#4-数据库迁移)
5. [灰度发布方案](#5-灰度发布方案)
6. [回滚操作](#6-回滚操作)
7. [监控与告警](#7-监控与告警)
8. [日常运维手册](#8-日常运维手册)

---

## 1. 前置条件

### 1.1 工具要求

| 工具 | 最低版本 | 用途 |
|------|---------|------|
| Docker | 24.0+ | 容器运行时 |
| Docker Compose | 2.20+ | 本地编排 |
| kubectl | 1.29+ | K8s 集群操作 |
| Go | 1.22+ | 后端构建 |
| Node.js | 20 LTS | 前端构建 |
| Python | 3.11+ | AI 引擎 |

### 1.2 服务器要求（生产）

| 角色 | 规格建议 | 说明 |
|------|---------|------|
| K8s Master | 4C/8G | 控制面 |
| K8s Worker (×3) | 8C/16G | 应用 Pod |
| 数据库节点 | 8C/32G + 500G SSD | PostgreSQL 主从 |
| Redis 节点 | 4C/8G | Redis 哨兵 |

### 1.3 必要密钥清单

在部署前准备以下密钥（通过 GitHub Secrets 配置）：

```
REGISTRY_USERNAME     # 镜像仓库用户名
REGISTRY_PASSWORD     # 镜像仓库密码
KUBECONFIG            # K8s 访问配置（base64 编码）
DB_PASSWORD           # PostgreSQL 密码
REDIS_PASSWORD        # Redis 密码
MINIO_ACCESS_KEY      # MinIO AccessKey
MINIO_SECRET_KEY      # MinIO SecretKey
MEILI_MASTER_KEY      # Meilisearch 主密钥
JWT_SECRET            # JWT 签名密钥（32 位随机字符串）
WECHAT_APP_ID         # 微信小程序 AppID
WECHAT_MCH_ID         # 微信商户号
WECHAT_API_KEY        # 微信支付 API 密钥
LLM_API_KEY           # LLM API 密钥
```

---

## 2. 开发环境快速启动

### 2.1 克隆与配置

```bash
git clone https://github.com/your-org/linximanager.git
cd linximanager

# 复制开发环境变量模板
cp deploy/env/.env.development deploy/env/.env.development.local
# 编辑 .env.development.local，填入本地所需配置
```

### 2.2 启动所有服务

```bash
# 一键启动全部依赖（PostgreSQL、Redis、MinIO、Meilisearch、AI 引擎、后端、管理台）
docker compose up -d

# 查看启动状态
docker compose ps

# 查看日志
docker compose logs -f backend
docker compose logs -f ai-engine
```

### 2.3 访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 管理后台 | http://localhost:3000 | Vue3 Admin |
| 后端 API | http://localhost:8080 | Go API |
| MinIO 控制台 | http://localhost:9001 | 对象存储管理 |
| Meilisearch | http://localhost:7700 | 搜索引擎 |
| AI 引擎 | http://localhost:8000/docs | FastAPI Swagger |

### 2.4 常用命令

```bash
# 重启单个服务
docker compose restart backend

# 重新构建并启动（代码变更后）
docker compose up -d --build backend

# 停止全部
docker compose down

# 停止并清除所有数据卷（慎用）
docker compose down -v
```

---

## 3. 生产环境部署

### 3.1 首次部署流程

#### Step 1：初始化 K8s 命名空间和密钥

```bash
# 创建命名空间
kubectl apply -f deploy/k8s/namespace.yaml

# 手动创建生产密钥（替换 CHANGE_ME 为 base64 编码的真实值）
# 生成 base64: echo -n "your_password" | base64
kubectl apply -f deploy/k8s/secret.yaml -n linximanager
```

#### Step 2：部署基础设施

```bash
# 部署 PostgreSQL（StatefulSet）
kubectl apply -f deploy/k8s/postgres-statefulset.yaml

# 部署 Redis
kubectl apply -f deploy/k8s/redis-deployment.yaml

# 等待数据库就绪
kubectl wait pod -l app=postgres -n linximanager --for=condition=Ready --timeout=120s
kubectl wait pod -l app=redis -n linximanager --for=condition=Ready --timeout=60s
```

#### Step 3：执行数据库初始化迁移

```bash
# 运行迁移 Job（见第 4 节详细说明）
kubectl apply -f deploy/k8s/migrate-job.yaml
kubectl wait job/db-migrate -n linximanager --for=condition=complete --timeout=300s
```

#### Step 4：部署应用服务

```bash
# 应用 ConfigMap
kubectl apply -f deploy/k8s/configmap.yaml

# 部署后端（含 HPA）
kubectl apply -f deploy/k8s/backend-deployment.yaml
kubectl apply -f deploy/k8s/backend-service.yaml

# 部署管理后台
kubectl apply -f deploy/k8s/admin-deployment.yaml
kubectl apply -f deploy/k8s/admin-service.yaml

# 配置 Ingress（需先安装 cert-manager 和 nginx-ingress-controller）
kubectl apply -f deploy/k8s/ingress.yaml
```

#### Step 5：验证部署

```bash
# 检查所有 Pod 状态
kubectl get pods -n linximanager

# 检查服务
kubectl get svc -n linximanager

# 检查 Ingress
kubectl get ingress -n linximanager

# 检查 HPA 状态
kubectl get hpa -n linximanager

# 查看后端日志
kubectl logs -l app=backend -n linximanager --tail=50

# 健康检查
curl https://linximanager.com/health
```

### 3.2 CI/CD 自动部署

推送到 `main` 分支会自动触发 GitHub Actions 部署流水线：

1. 构建 Docker 镜像（backend、admin、ai-engine）
2. 推送至镜像仓库（打 commit sha 标签）
3. 更新 K8s Deployment 镜像版本
4. 等待 Rolling Update 完成
5. 执行 smoke test

**流水线状态查看**：GitHub Actions → Deploy workflow

---

## 4. 数据库迁移

### 4.1 迁移工具

后端使用 GORM 的 `AutoMigrate` + 自定义迁移脚本，迁移文件位于 `backend/migrations/`。

### 4.2 执行迁移（开发环境）

```bash
# 进入后端容器执行迁移
docker compose exec backend /app/server migrate

# 或直接运行迁移命令
docker run --rm \
  --env-file deploy/env/.env.development.local \
  linximanager/backend:latest migrate
```

### 4.3 执行迁移（生产环境）

```bash
# 创建一次性迁移 Job
kubectl create job db-migrate-$(date +%Y%m%d) \
  --image=registry.cn-hangzhou.aliyuncs.com/linximanager/backend:latest \
  -n linximanager \
  -- /app/server migrate

# 查看迁移日志
kubectl logs job/db-migrate-$(date +%Y%m%d) -n linximanager
```

### 4.4 迁移注意事项

- 新增字段必须设置默认值或允许 NULL（避免锁表）
- 删除列/重命名列需分阶段：先保留旧列 → 代码上线 → 再删除旧列
- 大表 DDL 操作使用 `pg_repack` 或 `CREATE INDEX CONCURRENTLY` 避免长时间锁

---

## 5. 灰度发布方案

### 5.1 基于 K8s Rolling Update（默认）

生产环境默认使用 Rolling Update，每次替换 1 个 Pod：

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1        # 最多多出 1 个 Pod
    maxUnavailable: 0  # 始终保持全部可用
```

### 5.2 金丝雀发布（流量按比例切分）

适用于重大功能变更，需要 Nginx Ingress 支持：

```bash
# Step 1：部署 canary 版本
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend-canary
  namespace: linximanager
spec:
  replicas: 1  # 3 主 + 1 金丝雀 = 25% 流量
  selector:
    matchLabels:
      app: backend
      version: canary
  template:
    metadata:
      labels:
        app: backend
        version: canary
    spec:
      containers:
        - name: backend
          image: registry.cn-hangzhou.aliyuncs.com/linximanager/backend:canary-tag
EOF

# Step 2：观察错误率和延迟（至少 30 分钟）
kubectl logs -l app=backend,version=canary -n linximanager

# Step 3a：全量推进（金丝雀无异常）
kubectl set image deployment/backend backend=backend:canary-tag -n linximanager
kubectl delete deployment backend-canary -n linximanager

# Step 3b：回滚金丝雀（发现问题）
kubectl delete deployment backend-canary -n linximanager
```

---

## 6. 回滚操作

### 6.1 K8s 快速回滚

```bash
# 查看 Deployment 变更历史
kubectl rollout history deployment/backend -n linximanager

# 回滚到上一版本
kubectl rollout undo deployment/backend -n linximanager

# 回滚到指定版本
kubectl rollout undo deployment/backend --to-revision=3 -n linximanager

# 查看回滚状态
kubectl rollout status deployment/backend -n linximanager
```

### 6.2 数据库回滚

```bash
# 查看可用的迁移版本
kubectl exec -it pod/backend-xxx -n linximanager -- /app/server migrate status

# 回滚指定迁移
kubectl create job db-rollback-$(date +%Y%m%d) \
  --image=registry.cn-hangzhou.aliyuncs.com/linximanager/backend:current \
  -n linximanager \
  -- /app/server migrate down 1
```

### 6.3 完整版本回滚检查清单

- [ ] 确认目标回滚版本的镜像 tag
- [ ] 检查该版本的数据库 schema 兼容性
- [ ] 执行 `kubectl rollout undo` 回滚应用
- [ ] 如有 schema 不兼容，执行数据库回滚
- [ ] 验证健康检查端点：`curl https://linximanager.com/health`
- [ ] 观察错误率恢复正常（Grafana Dashboard）
- [ ] 通知团队并记录事故

---

## 7. 监控与告警

### 7.1 告警规则概览

| 告警名称 | 触发条件 | 严重级别 | 响应时限 |
|---------|---------|---------|---------|
| HighCPUUsage | CPU > 80% 持续 5min | warning | 30min |
| CriticalCPUUsage | CPU > 95% 持续 2min | critical | 立即 |
| HighMemoryUsage | 内存 > 85% 持续 5min | warning | 30min |
| DiskSpaceLow | 磁盘 > 85% 持续 10min | warning | 2h |
| HighErrorRate | 5xx > 1% 持续 2min | critical | 立即 |
| HighLatency | P95 > 500ms 持续 5min | warning | 15min |
| CriticalLatency | P99 > 2s 持续 3min | critical | 立即 |
| BackendPodDown | 可用副本 < 2 持续 1min | critical | 立即 |
| PostgresDown | 连接失败 持续 1min | critical | 立即 |
| RedisDown | 连接失败 持续 1min | critical | 立即 |
| RedisMemoryHigh | 内存 > 80% 持续 5min | warning | 30min |

### 7.2 Grafana Dashboard 导入

```bash
# 登录 Grafana 后：
# 左侧菜单 → Dashboards → Import → 上传文件
# 文件路径: deploy/monitoring/grafana-dashboard.json
```

### 7.3 Alertmanager 配置（企业微信通知示例）

```yaml
# alertmanager.yml
route:
  receiver: wechat-work
  group_by: [alertname, severity]
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h

receivers:
  - name: wechat-work
    webhook_configs:
      - url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY
        send_resolved: true
```

---

## 8. 日常运维手册

### 8.1 数据库备份

```bash
# 手动备份（在数据库节点执行）
kubectl exec -it pod/postgres-0 -n linximanager -- \
  pg_dump -U linxi linximanager | gzip > backup-$(date +%Y%m%d).sql.gz

# 上传至 MinIO
mc cp backup-$(date +%Y%m%d).sql.gz minio/linximanager-backups/postgres/

# 自动备份（推荐配置 CronJob）
kubectl apply -f - <<'EOF'
apiVersion: batch/v1
kind: CronJob
metadata:
  name: postgres-backup
  namespace: linximanager
spec:
  schedule: "0 2 * * *"  # 每天凌晨 2 点
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: backup
              image: postgres:16-alpine
              command:
                - sh
                - -c
                - |
                  pg_dump -h postgres-service -U $DB_USER $DB_NAME | \
                  gzip > /backup/backup-$(date +%Y%m%d).sql.gz
              env:
                - name: PGPASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: linxi-secret
                      key: DB_PASSWORD
          restartPolicy: OnFailure
EOF
```

### 8.2 后端扩容

```bash
# 手动扩容（HPA 自动扩容时不需要）
kubectl scale deployment backend --replicas=6 -n linximanager

# 调整 HPA 上限
kubectl patch hpa backend-hpa -n linximanager \
  --patch '{"spec":{"maxReplicas":20}}'
```

### 8.3 日志查看

```bash
# 实时日志（所有 backend Pod）
kubectl logs -l app=backend -n linximanager -f --max-log-requests=10

# 查看最近 100 行
kubectl logs deployment/backend -n linximanager --tail=100

# 指定时间范围（需要 stern 工具）
stern backend -n linximanager --since 1h

# 查看 AI 引擎日志
kubectl logs -l app=ai-engine -n linximanager -f
```

### 8.4 常见问题排查

**问题 1：后端 Pod 反复重启（CrashLoopBackOff）**

```bash
# 查看 Pod 详情
kubectl describe pod <pod-name> -n linximanager

# 查看最近一次崩溃的日志
kubectl logs <pod-name> -n linximanager --previous
```

常见原因：
- 数据库连接失败 → 检查 secret 中的 DB_PASSWORD 是否正确
- 内存 OOM → 调大 `resources.limits.memory`
- 健康检查超时 → 检查后端启动时间，调大 `initialDelaySeconds`

**问题 2：Ingress 返回 502**

```bash
# 检查后端 Pod 是否就绪
kubectl get pods -n linximanager -l app=backend

# 检查 Service 端点
kubectl get endpoints backend-service -n linximanager

# 查看 Ingress 控制器日志
kubectl logs -l app.kubernetes.io/name=ingress-nginx -n ingress-nginx --tail=50
```

**问题 3：数据库连接池耗尽**

```bash
# 查看当前连接数
kubectl exec -it pod/postgres-0 -n linximanager -- \
  psql -U linxi -d linximanager -c "SELECT count(*) FROM pg_stat_activity;"

# 查看连接分布
kubectl exec -it pod/postgres-0 -n linximanager -- \
  psql -U linxi -d linximanager -c \
  "SELECT client_addr, state, count(*) FROM pg_stat_activity GROUP BY 1,2 ORDER BY 3 DESC;"
```

**问题 4：Redis 内存告警**

```bash
# 查看内存使用
kubectl exec -it pod/redis-xxx -n linximanager -- redis-cli -a $REDIS_PASSWORD info memory

# 查看最大键
kubectl exec -it pod/redis-xxx -n linximanager -- \
  redis-cli -a $REDIS_PASSWORD --bigkeys

# 手动触发内存回收
kubectl exec -it pod/redis-xxx -n linximanager -- \
  redis-cli -a $REDIS_PASSWORD memory purge
```

### 8.5 证书续期

项目使用 cert-manager 自动续期 Let's Encrypt 证书，通常无需手动干预。

```bash
# 查看证书状态
kubectl get certificate -n linximanager

# 手动触发续期
kubectl annotate certificate linxi-tls-secret \
  cert-manager.io/issue-temporary-certificate="true" -n linximanager

# 查看 cert-manager 日志
kubectl logs -l app=cert-manager -n cert-manager --tail=50
```

### 8.6 性能基线参考

| 指标 | 正常范围 | 告警阈值 |
|------|---------|---------|
| API P95 延迟 | < 200ms | > 500ms |
| API P99 延迟 | < 500ms | > 2s |
| 5xx 错误率 | < 0.1% | > 1% |
| 后端 CPU | < 50% | > 70% |
| 后端内存 | < 512Mi/Pod | > 800Mi/Pod |
| PostgreSQL 连接数 | < 100 | > 160 |
| Redis 内存 | < 300Mi | > 400Mi |
