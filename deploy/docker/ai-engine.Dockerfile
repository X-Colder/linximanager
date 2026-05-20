# ===== Stage 1: Builder (依赖安装) =====
FROM python:3.11-slim AS builder

WORKDIR /build

RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir poetry

COPY pyproject.toml poetry.lock* requirements.txt* ./

RUN if [ -f pyproject.toml ]; then \
      poetry export -f requirements.txt --without-hashes -o /build/requirements_export.txt && \
      pip install --no-cache-dir --prefix=/install -r /build/requirements_export.txt; \
    elif [ -f requirements.txt ]; then \
      pip install --no-cache-dir --prefix=/install -r requirements.txt; \
    fi

# ===== Stage 2: Runner =====
FROM python:3.11-slim

RUN groupadd -r appgroup && \
    useradd -r -g appgroup -d /app -s /bin/false appuser

WORKDIR /app

COPY --from=builder /install /usr/local
COPY . .

RUN chown -R appuser:appgroup /app

USER appuser

ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8000/health')" || exit 1

CMD ["python", "-m", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "2"]
