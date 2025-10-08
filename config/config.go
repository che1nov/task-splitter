package config

import (
	"os"
	"strconv"

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
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "tasksplitter"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
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
