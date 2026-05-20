package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	RateLimit RateLimitConfig
}

type AppConfig struct {
	Env    string
	Port   string
	Secret string
}

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpire  time.Duration
	RefreshExpire time.Duration
}

type RateLimitConfig struct {
	QPS   float64
	Burst int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Env:    getEnv("APP_ENV", "development"),
			Port:   getEnv("APP_PORT", "8080"),
			Secret: getEnv("APP_SECRET", "changeme"),
		},
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			Name:            getEnv("DB_NAME", "linximanager"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME", 3600)) * time.Second,
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", "access-secret-key"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "refresh-secret-key"),
			AccessExpire:  time.Duration(getEnvInt("JWT_ACCESS_EXPIRE", 7200)) * time.Second,
			RefreshExpire: time.Duration(getEnvInt("JWT_REFRESH_EXPIRE", 604800)) * time.Second,
		},
		RateLimit: RateLimitConfig{
			QPS:   getEnvFloat("RATE_LIMIT_QPS", 100),
			Burst: getEnvInt("RATE_LIMIT_BURST", 200),
		},
	}
	return cfg, nil
}

func (c *DBConfig) DSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.Name +
		" sslmode=disable TimeZone=Asia/Shanghai"
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvFloat(key string, defaultVal float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return defaultVal
}
