"""AI 引擎配置，从环境变量读取"""

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    env: str = "development"
    port: int = 8000

    # 数据库
    database_url: str = "postgresql+asyncpg://linxi:linxi_dev_password@postgres:5432/linximanager"

    # Redis
    redis_url: str = "redis://:redis_dev_password@redis:6379/1"

    # LLM
    llm_api_key: str = ""
    llm_model: str = "claude-opus-4-7"

    class Config:
        env_file = ".env"
        extra = "ignore"


settings = Settings()
