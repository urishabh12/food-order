package config

import "os"

type Config struct {
	Addr     string
	BasePath string
}

func FromEnv() Config {
	return Config{
		Addr:     get("ADDR", ":8080"),
		BasePath: get("BASE_PATH", "/api"),
	}
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
