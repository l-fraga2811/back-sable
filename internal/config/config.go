package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Port           string
	SupabaseURL    string
	SupabaseKey    string
	JwksURL        string
	JwtSecret      string
	GoogleClientID string
	GoogleSecret   string
	DB             *gorm.DB
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

	// Initialize database connection
	db, err := NewGormDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return &Config{
		Port:           getEnv("PORT", ""),
		SupabaseURL:    getEnv("SUPABASE_URL", ""),
		SupabaseKey:    getEnv("SUPABASE_KEY", ""),
		JwksURL:        jwksURL,
		JwtSecret:      getEnv("SUPABASE_JWT_SECRET", ""),
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleSecret:   getEnv("GOOGLE_CLIENT_SECRET", ""),
		DB:             db,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func NewGormDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=America/Sao_Paulo",
		os.Getenv("SUPABASE_DB_HOST"),
		os.Getenv("SUPABASE_DB_USER"),
		os.Getenv("SUPABASE_DB_PASSWORD"),
		os.Getenv("SUPABASE_DB_NAME"),
		os.Getenv("SUPABASE_DB_PORT"),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
