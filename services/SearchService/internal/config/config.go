package config

import "os"

type (
	ElasticSearchURL string
	GRPCPort         string
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
