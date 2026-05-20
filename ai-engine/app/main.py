"""
灵犀掌柜 AI 引擎 - FastAPI 应用入口

服务职责：
- 销售预测（ARIMA / 机器学习模型）
- 智能补货建议
- 促销策略生成
- gRPC 接口（供 Go 后端调用）
- HTTP REST 接口（内部诊断/测试）
"""

from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.core.config import settings
from app.routers import health, predict


@asynccontextmanager
async def lifespan(app: FastAPI):
    """应用生命周期：启动时初始化资源，关闭时清理"""
    # 可在此初始化数据库连接池、加载模型等
    yield
    # 清理资源


app = FastAPI(
    title="灵犀掌柜 AI 引擎",
    version="1.0.0",
    description="销售预测与智能决策服务",
    docs_url="/docs" if settings.env == "development" else None,
    redoc_url=None,
    lifespan=lifespan,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # AI 引擎仅内网访问，不对外暴露
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(health.router)
app.include_router(predict.router, prefix="/api/v1")
