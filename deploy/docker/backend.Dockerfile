# ===== Stage 1: Builder =====
FROM m.daocloud.io/docker.io/library/golang:1.22-alpine3.19 AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -o /build/server ./cmd/server/main.go

# ===== Stage 2: Runner =====
FROM m.daocloud.io/docker.io/library/alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /build/server /app/server
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 数据库迁移文件
COPY --from=builder /build/migrations /app/migrations

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/server"]
