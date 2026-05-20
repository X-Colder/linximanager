# =============================================================================
# AI 引擎（Python FastAPI + gRPC）多阶段构建
#
# Build context: ./ai-engine
# =============================================================================

# ===== Stage 1: 依赖安装（利用层缓存）=====
FROM python:3.11-slim AS builder

WORKDIR /build

# 安装构建工具（仅 builder 阶段需要，不进入运行镜像）
RUN pip install --no-cache-dir --upgrade pip

# 优先复制依赖文件，使依赖层在代码未变化时可复用缓存
COPY requirements.txt ./

# 安装到独立前缀，方便复制到 runner 阶段
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt

# ===== Stage 2: 运行镜像 =====
FROM python:3.11-slim

# 创建非 root 用户运行应用
RUN groupadd -r appgroup && \
    useradd -r -g appgroup -d /app -s /bin/false appuser

WORKDIR /app

# 从 builder 阶段复制已安装的依赖
COPY --from=builder /install /usr/local

# 复制应用代码
COPY . .

RUN chown -R appuser:appgroup /app

USER appuser

# Python 最佳实践：不生成 .pyc 文件，日志直接输出到 stdout
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONPATH=/app

EXPOSE 8000
# gRPC 端口（由 gRPC server 在 app 内启动）
EXPOSE 50051

HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8000/health')" || exit 1

# workers=1：开发环境单进程，生产通过 --workers 参数或多副本扩展
CMD ["python", "-m", "uvicorn", "app.main:app", \
     "--host", "0.0.0.0", \
     "--port", "8000", \
     "--workers", "1"]
