package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	DBHost     string `koanf:"DB_HOST"`
	DBPort     int    `koanf:"DB_PORT"`
	DBUser     string `koanf:"DB_USER"`
	DBPassword string `koanf:"DB_PASS"`
	DBName     string `koanf:"DB_NAME"`

	AppPort string `koanf:"APP_PORT"`
	AppEnv  string `koanf:"APP_ENV"`
	AppHost string `koanf:"APP_HOST"`

	JWTSecret string `koanf:"JWT_SECRET"`

	MinIOEndpoint  string `koanf:"MINIO_ENDPOINT"`
	MinIOAccessKey string `koanf:"MINIO_ACCESS_KEY"`
	MinIOSecretKey string `koanf:"MINIO_SECRET_KEY"`
	MinIOBucket    string `koanf:"MINIO_BUCKET"`
	MinIOUseSSL    string `koanf:"MINIO_USE_SSL"`

	RedisAddr     string `koanf:"REDIS_ADDR"`
	RedisPassword string `koanf:"REDIS_PASSWORD"`
	GinMode       string `koanf:"GIN_MODE"`

	// CORS Configuration
	CORSAllowedOrigins   string `koanf:"CORS_ALLOWED_ORIGINS"`
	CORSAllowedMethods   string `koanf:"CORS_ALLOWED_METHODS"`
	CORSAllowedHeaders   string `koanf:"CORS_ALLOWED_HEADERS"`
	CORSExposeHeaders    string `koanf:"CORS_EXPOSE_HEADERS"`
	CORSAllowCredentials bool   `koanf:"CORS_ALLOW_CREDENTIALS"`
	CORSMaxAge           int    `koanf:"CORS_MAX_AGE"`

	// Rate Limit
	RateLimitEnabled   bool `koanf:"RATE_LIMIT_ENABLED"`
	RateLimitPerSecond int  `koanf:"RATE_LIMIT_PER_SECOND"`
	RateLimitBurst     int  `koanf:"RATE_LIMIT_BURST"`
}

func LoadConfig() (*Config, error) {
	k := koanf.New(".")

	configPath := os.Getenv("CONFIG_FILE_PATH")
	if configPath == "" {
		configPath = ".env"
	}

	if err := k.Load(file.Provider(configPath), dotenv.Parser()); err != nil {
		if os.IsNotExist(err) {
			log.Default().Printf("Info: config file '%s' not found, using environment variables only", configPath)
		} else {
			log.Default().Printf("Warning: error reading config file '%s': %v", configPath, err)
		}
	}

	if err := k.Load(env.Provider("", ".", nil), nil); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	config := &Config{}
	if err := k.Unmarshal("", config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort)
}

func (c *Config) IsProduction() bool {
	return strings.ToLower(c.AppEnv) == "production" || strings.ToLower(c.AppEnv) == "release"
}

func (c *Config) GetCORSAllowedOrigins() []string {
	if c.CORSAllowedOrigins == "" {
		return []string{"http://localhost:3000", "http://localhost:3001", "https://yourdomain.com"}
	}
	return strings.Split(c.CORSAllowedOrigins, ",")
}

func (c *Config) GetCORSAllowedMethods() []string {
	if c.CORSAllowedMethods == "" {
		return []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	return strings.Split(c.CORSAllowedMethods, ",")
}

func (c *Config) GetCORSAllowedHeaders() []string {
	if c.CORSAllowedHeaders == "" {
		return []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	}
	return strings.Split(c.CORSAllowedHeaders, ",")
}

func (c *Config) GetCORSExposeHeaders() []string {
	if c.CORSExposeHeaders == "" {
		return []string{"Content-Length"}
	}
	return strings.Split(c.CORSExposeHeaders, ",")
}

func (c *Config) GetCORSMaxAge() time.Duration {
	if c.CORSMaxAge == 0 {
		return 12 * time.Hour
	}
	return time.Duration(c.CORSMaxAge) * time.Hour
}
