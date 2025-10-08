package main

import (
	"fmt"
	"os"
	"task-splitter/config"
)

func main() {
	fmt.Println("🧪 Testing DATABASE_URL parsing...")

	// Тестируем с реальным Railway DATABASE_URL
	testURL := "postgres://postgres:password123@containers-us-west-123.railway.app:6543/railway?sslmode=require"
	os.Setenv("DATABASE_URL", testURL)

	cfg := config.Load()
	dbConfig := cfg.Database

	fmt.Printf("✅ Parsed config:\n")
	fmt.Printf("   Host: %s\n", dbConfig.Host)
	fmt.Printf("   Port: %s\n", dbConfig.Port)
	fmt.Printf("   User: %s\n", dbConfig.User)
	fmt.Printf("   Password: %s\n", "***")
	fmt.Printf("   DBName: %s\n", dbConfig.DBName)
	fmt.Printf("   SSLMode: %s\n", dbConfig.SSLMode)

	// Проверяем, что парсинг работает правильно
	expectedHost := "containers-us-west-123.railway.app"
	expectedPort := "6543"
	expectedUser := "postgres"
	expectedDBName := "railway"
	expectedSSLMode := "require"

	if dbConfig.Host != expectedHost {
		fmt.Printf("❌ Host mismatch: got %s, expected %s\n", dbConfig.Host, expectedHost)
	} else {
		fmt.Printf("✅ Host correct: %s\n", dbConfig.Host)
	}

	if dbConfig.Port != expectedPort {
		fmt.Printf("❌ Port mismatch: got %s, expected %s\n", dbConfig.Port, expectedPort)
	} else {
		fmt.Printf("✅ Port correct: %s\n", dbConfig.Port)
	}

	if dbConfig.User != expectedUser {
		fmt.Printf("❌ User mismatch: got %s, expected %s\n", dbConfig.User, expectedUser)
	} else {
		fmt.Printf("✅ User correct: %s\n", dbConfig.User)
	}

	if dbConfig.DBName != expectedDBName {
		fmt.Printf("❌ DBName mismatch: got %s, expected %s\n", dbConfig.DBName, expectedDBName)
	} else {
		fmt.Printf("✅ DBName correct: %s\n", dbConfig.DBName)
	}

	if dbConfig.SSLMode != expectedSSLMode {
		fmt.Printf("❌ SSLMode mismatch: got %s, expected %s\n", dbConfig.SSLMode, expectedSSLMode)
	} else {
		fmt.Printf("✅ SSLMode correct: %s\n", dbConfig.SSLMode)
	}

	fmt.Println("\n🎯 Test completed!")
}
