package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASS"`
	DBName     string `mapstructure:"DB_NAME"`

	AppPort string `mapstructure:"APP_PORT"`
	AppEnv  string `mapstructure:"APP_ENV"`
	AppHost string `mapstructure:"APP_HOST"`

	JWTSecret string `mapstructure:"JWT_SECRET"`

	MinIOEndpoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinIOAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinIOSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinIOBucket    string `mapstructure:"MINIO_BUCKET"`
	MinIOUseSSL    string `mapstructure:"MINIO_USE_SSL"`

	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	GinMode       string `mapstructure:"GIN_MODE"`

	// CORS Configuration
	CORSAllowedOrigins   string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CORSAllowedMethods   string `mapstructure:"CORS_ALLOWED_METHODS"`
	CORSAllowedHeaders   string `mapstructure:"CORS_ALLOWED_HEADERS"`
	CORSExposeHeaders    string `mapstructure:"CORS_EXPOSE_HEADERS"`
	CORSAllowCredentials bool   `mapstructure:"CORS_ALLOW_CREDENTIALS"`
	CORSMaxAge           int    `mapstructure:"CORS_MAX_AGE"`

	// Rate Limit
	RateLimitEnabled   bool `mapstructure:"RATE_LIMIT_ENABLED"`
	RateLimitPerSecond int  `mapstructure:"RATE_LIMIT_PER_SECOND"`
	RateLimitBurst     int  `mapstructure:"RATE_LIMIT_BURST"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigType("env")
	v.AutomaticEnv()

	configPath := os.Getenv("CONFIG_FILE_PATH")
	if configPath == "" {
		configPath = ".env"
	}

	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", configPath, err)
	}

	config := &Config{}

	if err := v.Unmarshal(config); err != nil {
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
