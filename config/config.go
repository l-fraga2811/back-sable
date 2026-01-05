package config

import (
    "log"
    "os"
    "gorm.io/gorm"
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
    DB             *gorm.DB  // Nova propriedade
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

    checkPort := getEnv("PORT", "")
    if checkPort == "" {
        log.Fatal("PORT environment variable is not set")
    }

    // Inicializar GORM
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