package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	SupabaseURL    string
	SupabaseKey    string
	JwksURL        string
	JwtSecret      string
	GoogleClientID string
	GoogleSecret   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	supabaseURL := getEnv("SUPABASE_URL", "")
	jwksURL := getEnv("SUPABASE_JWKS_URL", "")
	if jwksURL == "" && supabaseURL != "" {
		jwksURL = strings.TrimRight(supabaseURL, "/") + "/auth/v1/.well-known/jwks.json"
	}

	return &Config{
		Port:           getEnv("PORT", "3000"),
		SupabaseURL:    supabaseURL,
		SupabaseKey:    getEnv("SUPABASE_KEY", ""),
		JwksURL:        jwksURL,
		JwtSecret:      getEnv("SUPABASE_JWT_SECRET", ""),
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleSecret:   getEnv("GOOGLE_CLIENT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
