"""健康检查路由"""

from fastapi import APIRouter

router = APIRouter()


@router.get("/health", tags=["health"])
async def health():
    """Liveness 探针，返回 200 即表示服务正常运行"""
    return {"status": "ok"}


@router.get("/ready", tags=["health"])
async def ready():
    """Readiness 探针，可在此加入数据库/Redis 连通性检查"""
    return {"status": "ready"}
