package config

import (
    "fmt"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

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