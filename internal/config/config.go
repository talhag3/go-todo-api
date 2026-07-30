package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser  string
	DBPass  string
	DBName  string
	AppPort int
	AppName string
}

func LoadConf() (*Config, error) {

	err := godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("error loading the env file: %w", err)
	}

	return &Config{
		DBUser:  getEnv("DB_USER", "username"),
		DBPass:  getEnv("DB_PASS", ""),
		DBName:  getEnv("DB_NAME", ""),
		AppPort: getEnvAsInt("APP_PORT", 8080),
		AppName: getEnv("APP_NAME", "myapp"),
	}, nil
}

/*
@enVar string # env variable
@dVal string # default value
*/
func getEnv(enVar string, dVal string) string {
	// Get the Env Value for the env
	value, ok := os.LookupEnv(enVar)
	if ok {
		return value
	}
	return dVal
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func ConnectionString(config *Config) string {
	return fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		config.DBUser, config.DBPass, config.DBName)
}
