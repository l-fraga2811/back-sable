package config

import (
	"log"
	"os"

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


	jwksURL := getEnv("SUPABASE_JWKS_URL", "")
	if jwksURL == "" {
		jwksURL = getEnv("NEXT_PUBLIC_SUPABASE_JWKS_URL", "")
	}

	return &Config{
		Port:           getEnv("PORT", "3000"),
		SupabaseURL:    getEnv("SUPABASE_URL", ""),
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
