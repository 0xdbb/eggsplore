package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the configuration values from environment variables.
type Config struct {
	Production string
	Port       string

	DbUrl string

	RMQAddress    string
	RMQAddressDev string

	RedisDB       int
	RedisAddr     string
	RedisPassword string

	Recipients string
	AdminEmail string

	S3Bucket  string
	AWSRegion string

	AppDomain string

	ArkeselSMSApiKey string
	ArkselSenderID   string

	CitizenReportAppName string

	ResendApiKey string

	TokenSecret                  string
	AccessTokenDuration          time.Duration
	AccountApprovalTokenDuration time.Duration
	RefreshTokenDuration         time.Duration

	// WS
	WSTableName string
	WSEndpoint  string
}

// LoadConfig loads environment variables from the .env file (if it exists)
// or falls back to system environment variables (as in AWS App Runner).
func LoadConfig(path ...string) (*Config, error) {
	// Default .env path
	envFile := ".env"
	if len(path) > 0 {
		envFile = path[0]
	}

	// Load .env but don’t crash if missing (useful in AWS)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("⚠️  No %s file found, relying on system environment variables", envFile)
	}

	// Parse durations safely
	accessTokenDuration, err := parseDuration("ACCESS_TOKEN_DURATION")
	if err != nil {
		return nil, err
	}
	refreshTokenDuration, err := parseDuration("REFRESH_TOKEN_DURATION")
	if err != nil {
		return nil, err
	}
	accountApprovalTokenDuration, err := parseDuration("ACCOUNT_APPROVAL_TOKEN_DURATION")
	if err != nil {
		return nil, err
	}

	// Populate config
	config := &Config{
		Production: os.Getenv("PRODUCTION"),
		Port:       os.Getenv("PORT"),

		DbUrl: os.Getenv("DB_URL"),

		ResendApiKey: os.Getenv("RESEND_API_KEY"),

		CitizenReportAppName: os.Getenv("APP_NAME"),
		AppDomain:            os.Getenv("APP_DOMAIN"),

		ArkeselSMSApiKey: os.Getenv("ARKESEL_API_KEY"),
		ArkselSenderID:   os.Getenv("ARKESEL_SENDER_ID"),

		RedisAddr:     os.Getenv("REDIS_URL"),
		RedisDB:       0,
		RedisPassword: os.Getenv("REDIS_PASSWORD"),

		TokenSecret:                  os.Getenv("TOKEN_SECRET"),
		AccountApprovalTokenDuration: accountApprovalTokenDuration,
		AccessTokenDuration:          accessTokenDuration,
		RefreshTokenDuration:         refreshTokenDuration,

		Recipients: os.Getenv("RECIPIENTS"),
		AdminEmail: os.Getenv("ADMIN_EMAIL"),

		RMQAddress:    os.Getenv("RMQ_ADDRESS"),
		RMQAddressDev: os.Getenv("RMQ_ADDRESS_DEV"),

		S3Bucket:  os.Getenv("S3_BUCKET"),
		AWSRegion: os.Getenv("AWS_REGION"),

		WSEndpoint:  os.Getenv("WS_ENDPOINT"),
		WSTableName: os.Getenv("WS_TABLE_NAME"),
	}

	// Validate required vars
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// parseDuration pulls an env var and parses it into a time.Duration.
func parseDuration(envKey string) (time.Duration, error) {
	val := os.Getenv(envKey)
	if val == "" {
		return 0, fmt.Errorf("missing required duration env var: %s", envKey)
	}
	dur, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s: %w", envKey, err)
	}
	return dur, nil
}

// validateConfig checks that all required environment variables are set.
func validateConfig(config *Config) error {
	if config.DbUrl == "" {
		return errors.New("missing required environment variable: DB_URL")
	}
	if config.AccessTokenDuration == 0 {
		return errors.New("missing or invalid required environment variable: ACCESS_TOKEN_DURATION")
	}
	if config.RefreshTokenDuration == 0 {
		return errors.New("missing or invalid required environment variable: REFRESH_TOKEN_DURATION")
	}
	if config.ResendApiKey == "" {
		return errors.New("missing required environment variable: RESEND_API_KEY")
	}
	if config.ArkeselSMSApiKey == "" {
		return errors.New("missing required environment variable: ARKESEL_API_KEY")
	}
	if config.AccountApprovalTokenDuration == 0 {
		return errors.New("missing or invalid required environment variable: ACCOUNT_APPROVAL_TOKEN_DURATION")
	}
	if config.RMQAddress == "" {
		return errors.New("missing or invalid required environment variable: RMQ_ADDRESS")
	}
	if config.Port == "" {
		return errors.New("missing required environment variable: PORT")
	}
	return nil
}
