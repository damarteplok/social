package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr         string
	DbAddr       string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
	ENV          string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return Config{
		Addr:         GetString("ADDR", ":8080"),
		DbAddr:       GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable"),
		MaxOpenConns: GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns: GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  GetString("DB_MAX_IDLE_TIME", "15m"),
		ENV:          GetString("ENV", "development"),
	}
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
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
