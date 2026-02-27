package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Kafka    KafkaConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string
	Postgres PostgresConfig
	Oracle   OracleConfig
}

// PostgresConfig holds PostgreSQL specific configuration
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// OracleConfig holds Oracle specific configuration
type OracleConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	ServiceName string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// JWT
type JWTConfig struct {
	SecretKey string
	ExpiresIn int
}

// KafkaConfig holds Kafka configuration.
// Brokers supports single broker or cluster:
//   - Single : KAFKA_BROKERS=localhost:9092
//   - Cluster: KAFKA_BROKERS=broker1:9092,broker2:9092,broker3:9092
type KafkaConfig struct {
	Brokers          []string
	GroupID          string
	AutoOffsetReset  string
	SecurityProtocol string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	config := &Config{
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "postgres"),
			Postgres: PostgresConfig{
				Host:     getEnv("POSTGRES_HOST", "localhost"),
				Port:     getEnv("POSTGRES_PORT", "5432"),
				User:     getEnv("POSTGRES_USER", "postgres"),
				Password: getEnv("POSTGRES_PASSWORD", "postgres"),
				DBName:   getEnv("POSTGRES_DB", "hyperscal_db"),
			},
			Oracle: OracleConfig{
				Host:        getEnv("ORACLE_HOST", "localhost"),
				Port:        getEnv("ORACLE_PORT", "1521"),
				User:        getEnv("ORACLE_USER", "system"),
				Password:    getEnv("ORACLE_PASSWORD", "oracle"),
				ServiceName: getEnv("ORACLE_SERVICE_NAME", "ORCLPDB1"),
			},
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", "MySecretKey123"),
			ExpiresIn: getEnvInt("JWT_EXPIRES_IN", 24),
		},
		Kafka: KafkaConfig{
			Brokers:          getEnvSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			GroupID:          getEnv("KAFKA_GROUP_ID", "hyperscal-consumer-group"),
			AutoOffsetReset:  getEnv("KAFKA_AUTO_OFFSET_RESET", "earliest"),
			SecurityProtocol: getEnv("KAFKA_SECURITY_PROTOCOL", "PLAINTEXT"),
		},
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		// Parse string ke int
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvSlice gets a comma-separated environment variable as a string slice.
// Example: KAFKA_BROKERS=broker1:9092,broker2:9092 → ["broker1:9092", "broker2:9092"]
func getEnvSlice(key string, defaultVal []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
