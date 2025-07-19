package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	Port        int
	Environment string
	
	// Database configuration
	DatabaseURL string
	
	// Kubernetes configuration
	KubeConfigPath string
	Namespace      string
	MockKubernetes bool
	
	// Monitoring configuration
	PrometheusEnabled bool
	GrafanaURL        string
}

// init loads environment variables from .env file if it exists
func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Using environment variables.")
	}
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnvOrDefault("PORT", "8080"))
	if err != nil {
		port = 8080
	}
	
	return &Config{
		Port:              port,
		Environment:       getEnvOrDefault("ENVIRONMENT", "development"),
		DatabaseURL:       getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/chaos_platform?sslmode=disable"),
		KubeConfigPath:    getEnvOrDefault("KUBECONFIG", ""),
		Namespace:         getEnvOrDefault("NAMESPACE", "default"),
		MockKubernetes:    getEnvOrDefault("MOCK_KUBERNETES", "false") == "true",
		PrometheusEnabled: getEnvOrDefault("PROMETHEUS_ENABLED", "true") == "true",
		GrafanaURL:        getEnvOrDefault("GRAFANA_URL", "http://localhost:3000"),
	}, nil
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}