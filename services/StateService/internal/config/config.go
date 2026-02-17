package config

import "os"

type (
	Host      string
	RedisPort string
	GRPCPort  string
	Password  string
	DB        int
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
