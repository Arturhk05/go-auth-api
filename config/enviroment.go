package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
	DbSSLMode  string
}

var Env *Environment

func init() {
	Env = NewEnvironment()
}

func NewEnvironment() *Environment {
	godotenv.Load()

	return &Environment{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     getEnvironmentVariableAsInt("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		DbSSLMode:  os.Getenv("DB_SSL_MODE"),
	}
}

func getEnvironmentVariableAsInt(key string) int {
	value := os.Getenv(key)

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return intValue
}
