"""AI 引擎配置"""

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    env: str = "development"
    port: int = 8000

    database_url: str = "postgresql+asyncpg://linxi:linxi_dev_password@postgres:5432/linximanager"
    redis_url: str = "redis://:redis_dev_password@redis:6379/1"

    llm_api_key: str = ""
    llm_model: str = "claude-opus-4-7"

    class Config:
        env_file = ".env"
        extra = "ignore"


settings = Settings()
