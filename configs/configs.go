package configs

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

func getEnvString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	var intVal int
	_, err := fmt.Sscanf(val, "%d", &intVal)
	if err != nil {
		return fallback
	}
	return intVal
}

func getEnvBytes(key string, fallback []byte) []byte {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	decoded, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return fallback
	}
	return decoded
}

var (
	// Configuration constants
	GIN_MODE         = getEnvString("GIN_MODE", gin.DebugMode)
	GORM_LOGGER_MODE = logger.LogLevel(getEnvInt("GORM_LOGGER_MODE", int(logger.Info)))
	HTTP_PORT        = getEnvInt("HTTP_PORT", 8080)
	CURRENT_DOMAIN   = getEnvString("CURRENT_DOMAIN", "localhost")

	// Database configuration
	DB_TYPE     = getEnvString("DB_TYPE", "mysql")
	DB_USER     = getEnvString("DB_USER", "username")
	DB_PASSWORD = getEnvString("DB_PASSWORD", "password")
	DB_HOST     = getEnvString("DB_HOST", "127.0.0.1:3306")
	DB_NAME     = getEnvString("DB_NAME", "dbname")

	// Gmail SMTP configuration
	GMAIL_ACCOUNT       = getEnvString("GMAIL_ACCOUNT", "")
	GMAIL_CLIENT_ID     = getEnvString("GMAIL_CLIENT_ID", "")
	GMAIL_CLIENT_SECRET = getEnvString("GMAIL_CLIENT_SECRET", "")
	GMAIL_REFRESH_TOKEN = getEnvString("GMAIL_REFRESH_TOKEN", "")

	// CAPTCHA configuration
	CAPTCHA_SECRET_KEY = getEnvString("CAPTCHA_SECRET_KEY", "")

	// User password salt
	USER_PSW_SALT = getEnvString("USER_PSW_SALT", "a1b2c3d4e5f6g7h8i9j0")

	// FB token encryption keys
	FBTOKEN_KEY = getEnvBytes("FBTOKEN_KEY", []byte{
		0x6b, 0x65, 0x79, 0x31, 0x32, 0x33, 0x34, 0x35,
		0x5a, 0x58, 0x43, 0x56, 0x42, 0x4e, 0x4d, 0x2c,
		0x2e, 0x2f, 0x3b, 0x3a, 0x27, 0x22, 0x5f, 0x2d,
		0x3d, 0x2b, 0x5b, 0x5d, 0x7b, 0x7d, 0x3c, 0x3e,
	})
	FBTOKEN_IV = getEnvBytes("FBTOKEN_IV", []byte{
		0x22, 0x87, 0x70, 0xf5, 0x13, 0xa4, 0xad, 0xdd,
		0x5e, 0x6a, 0x5f, 0x4d, 0x02, 0xae, 0x80, 0x69,
	})
)
