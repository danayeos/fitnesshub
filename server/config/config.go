// config/config.go
package config

import (
	"os"
	"strings"
)

// Конфигурация для подключения к базе данных
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var AppConfig Config

// Загрузка конфигурации из переменных окружения
func LoadConfig() {
	// Проверяем, есть ли DATABASE_URL в переменных окружения (Render его использует)
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL != "" {
		AppConfig = Config{
			DBHost:     databaseURL, // Сохраняем полный URL
			DBPort:     "",
			DBUser:     "",
			DBPassword: "",
			DBName:     "",
		}
	} else {
		// Если переменной нет, используем локальные настройки
		AppConfig = Config{
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", "123456"),
			DBName:     getEnv("DB_NAME", "mydb"),
		}
	}
}

// Получение значения из переменной окружения или использование значения по умолчанию
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
