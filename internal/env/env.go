package env

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr                 string
	DbAddr               string
	MaxOpenConns         int
	MaxIdleConns         int
	MaxIdleTime          string
	ENV                  string
	ApiUrl               string
	ZeebeAddr            string
	ZeebeClientID        string
	ZeebeClientSecret    string
	ZeebeAuthServerUrl   string
	FrontendURL          string
	MailerFromEmail      string
	MailerApiKey         string
	MailerExp            time.Duration
	AdminUser            string
	AdminPass            string
	JwtSecret            string
	JwtIss               string
	JwtExp               time.Duration
	RedisAddr            string
	RedisPass            string
	RedisDB              int
	RedisEnabled         bool
	AllowedOrigin        []string
	RequestPerTimeFrame  int
	RateLimiterEnabled   bool
	RateLimiterTimeFrame time.Duration
	MinioEndPoint        string
	MinioPort            int
	MinioSSL             bool
	MinioAccessKey       string
	MinioSecretKey       string
	MinioDefaultBucket   string
	MinioExpires         time.Duration
	MinioEnabled         bool
}

var Envs = initConfig()

func initConfig() Config {
	if os.Getenv("ENV") != "test" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, continuing without it.")
		}
	}

	return Config{
		Addr:                 GetString("ADDR", ":8080"),
		DbAddr:               GetString("DB_ADDR", ""),
		MaxOpenConns:         GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns:         GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:          GetString("DB_MAX_IDLE_TIME", "15m"),
		ENV:                  GetString("ENV", "development"),
		FrontendURL:          GetString("FRONTEND_URL", "http://localhost:5173"),
		MailerFromEmail:      GetString("MAILIER_FROM_EMAIL", "damar@test.com"),
		MailerApiKey:         GetString("MAILIER_API_KEY", ""),
		MailerExp:            GetDay("MAILER_EXP", 3),
		AdminUser:            GetString("ADMIN_USER", "admin"),
		AdminPass:            GetString("ADMIN_PASS", "admin"),
		JwtSecret:            GetString("JWT_SECRET", "admin"),
		JwtIss:               GetString("JWT_ISS", "damar"),
		JwtExp:               GetDay("JWT_EXP", 3),
		RedisAddr:            GetString("REDIS_ADDR", "localhost:6379"),
		RedisPass:            GetString("REDIS_PASS", ""),
		RedisDB:              GetInt("REDIS_DB", 0),
		RedisEnabled:         GetBool("REDIS_ENABLED", false),
		AllowedOrigin:        GetStringSlice("CORS_ALLOWED_ORIGIN", "https://*,http://*"),
		ZeebeAddr:            GetString("ZEEBE_ADDR", "localhost:8080"),
		ZeebeClientID:        GetString("ZEEBE_CLIENT_ID", "localhost:8080"),
		ZeebeClientSecret:    GetString("ZEEBE_CLIENT_SECRET", "localhost:8080"),
		ZeebeAuthServerUrl:   GetString("ZEEBE_AUTH_SERVER_URL", "localhost:8080"),
		RequestPerTimeFrame:  GetInt("REQUEST_PER_TIME_FRAME", 60),
		RateLimiterEnabled:   GetBool("RATE_LIMITER_ENABLED", true),
		RateLimiterTimeFrame: GetTimeSecond("RATE_LIMITER_TIME_FRAME", 5),
		MinioEndPoint:        GetString("MINIO_ENDPOINT", "127.0.0.1"),
		MinioPort:            GetInt("MINIO_PORT", 9000),
		MinioSSL:             GetBool("MINIO_SSL", false),
		MinioAccessKey:       GetString("MINIO_ACCESS_KEY", ""),
		MinioSecretKey:       GetString("MINIO_SECRET_KEY", ""),
		MinioDefaultBucket:   GetString("MINIO_DEFAULT_BUCKET", ""),
		MinioExpires:         GetDay("MINIO_EXPIRES", 1),
		MinioEnabled:         GetBool("MINIO_ENABLED", true),
	}
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func GetStringSlice(key, fallback string) []string {
	val := GetString(key, fallback)
	return strings.Split(val, ",")
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsBool, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return valAsBool
}

func GetTimeSecond(key string, fallback int) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return time.Second * time.Duration(fallback)
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return time.Second * time.Duration(fallback)
	}

	return time.Second * time.Duration(valAsInt)
}

func GetDay(key string, fallback int) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return time.Hour * 24 * time.Duration(fallback)
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return time.Hour * 24 * time.Duration(fallback)
	}

	return time.Hour * 24 * time.Duration(valAsInt)
}
