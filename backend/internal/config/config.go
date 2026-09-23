// Package config 负责加载应用配置。
package config

import (
	"log"
	"os"
	"strconv"
)

// Config 应用配置。
type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPassword    string
	JWTSecret     string
	JWTExpireSec  int
	LogLevel      string
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

// Load 从环境变量加载配置，缺失时使用默认值（本地开发）。
func Load() *Config {
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "127.0.0.1"),
		DBPort:        getEnv("DB_PORT", "44005"),
		DBName:        getEnv("DB_NAME", "esportsbar_db"),
		DBUser:        getEnv("DB_USER", "esportsbar_user"),
		DBPassword:    getEnv("DB_PASSWORD", "esportsbar_pwd"),
		JWTSecret:     getEnv("JWT_SECRET", "change_me_to_a_long_random_string"),
		JWTExpireSec:  getEnvInt("JWT_EXPIRE_SEC", 86400),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		RedisHost:     getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:     getEnv("REDIS_PORT", "46305"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}
}

// DSN 返回 GORM MySQL 连接串。
func (c *Config) DSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName +
		"?charset=utf8mb4&parseTime=True&loc=Local"
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Printf("invalid int env %s=%s, use default %d", key, v, def)
	}
	return def
}
