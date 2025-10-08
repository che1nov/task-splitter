package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	Keycloak KeycloakConfig
	OpenAI   OpenAIConfig
	GigaChat GigaChatConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	VHost    string
}

type KeycloakConfig struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
	AdminUser    string
	AdminPass    string
}

type OpenAIConfig struct {
	APIKey string
	Model  string
}

type GigaChatConfig struct {
	ClientID string
	Scope    string
	AuthKey  string
}

func Load() *Config {
	// Загружаем .env файл если он существует
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: loadDatabaseConfig(),
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
			User:     getEnv("RABBITMQ_USER", "guest"),
			Password: getEnv("RABBITMQ_PASSWORD", "guest"),
			VHost:    getEnv("RABBITMQ_VHOST", "/"),
		},
		Keycloak: KeycloakConfig{
			URL:          getEnv("KEYCLOAK_URL", "http://localhost:8081"),
			Realm:        getEnv("KEYCLOAK_REALM", "tasksplitter"),
			ClientID:     getEnv("KEYCLOAK_CLIENT_ID", "tasksplitter-api"),
			ClientSecret: getEnv("KEYCLOAK_CLIENT_SECRET", ""),
			AdminUser:    getEnv("KEYCLOAK_ADMIN_USER", "admin"),
			AdminPass:    getEnv("KEYCLOAK_ADMIN_PASS", "admin"),
		},
		OpenAI: OpenAIConfig{
			APIKey: getEnv("OPENAI_API_KEY", ""),
			Model:  getEnv("OPENAI_MODEL", "gpt-4"),
		},
		GigaChat: GigaChatConfig{
			ClientID: getEnv("GIGACHAT_CLIENT_ID", ""),
			Scope:    getEnv("GIGACHAT_SCOPE", "GIGACHAT_API_PERS"),
			AuthKey:  getEnv("GIGACHAT_AUTH_KEY", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// loadDatabaseConfig загружает конфигурацию базы данных с поддержкой DATABASE_URL
func loadDatabaseConfig() DatabaseConfig {
	// Проверяем наличие DATABASE_URL (Railway предоставляет это)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		fmt.Printf("🔗 Found DATABASE_URL: %s\n", databaseURL)
		// Парсим DATABASE_URL
		// Формат: postgres://user:password@host:port/database?sslmode=require
		config := parseDatabaseURL(databaseURL)
		fmt.Printf("📊 Parsed DB config: host=%s, port=%s, user=%s, dbname=%s, sslmode=%s\n",
			config.Host, config.Port, config.User, config.DBName, config.SSLMode)
		return config
	}

	fmt.Println("⚠️  DATABASE_URL not found, using individual variables")
	// Используем отдельные переменные
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "tasksplitter"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// parseDatabaseURL парсит DATABASE_URL в DatabaseConfig
func parseDatabaseURL(databaseURL string) DatabaseConfig {
	// Простой парсинг DATABASE_URL
	// Формат: postgres://user:password@host:port/database?sslmode=require

	// Убираем префикс postgres://
	url := strings.TrimPrefix(databaseURL, "postgres://")

	// Разделяем на части
	parts := strings.Split(url, "@")
	if len(parts) != 2 {
		// Если не удалось распарсить, используем значения по умолчанию
		return DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "tasksplitter",
			SSLMode:  "require",
		}
	}

	// Парсим user:password
	userPass := strings.Split(parts[0], ":")
	user := "postgres"
	password := "postgres"
	if len(userPass) == 2 {
		user = userPass[0]
		password = userPass[1]
	}

	// Парсим host:port/database
	hostPortDB := strings.Split(parts[1], "/")
	if len(hostPortDB) != 2 {
		return DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     user,
			Password: password,
			DBName:   "tasksplitter",
			SSLMode:  "require",
		}
	}

	// Парсим host:port
	hostPort := strings.Split(hostPortDB[0], ":")
	host := "localhost"
	port := "5432"
	if len(hostPort) == 2 {
		host = hostPort[0]
		port = hostPort[1]
	}

	// Парсим database и параметры
	dbName := hostPortDB[1]
	sslMode := "require"

	// Убираем параметры из имени базы данных
	if strings.Contains(dbName, "?") {
		dbParts := strings.Split(dbName, "?")
		dbName = dbParts[0]
		// Проверяем sslmode в параметрах
		if strings.Contains(dbParts[1], "sslmode=") {
			sslParts := strings.Split(dbParts[1], "sslmode=")
			if len(sslParts) > 1 {
				sslMode = strings.Split(sslParts[1], "&")[0]
			}
		}
	}

	return DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}
}
