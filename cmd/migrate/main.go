// cmd/migrate/main.go
package main

import (
    "log"
    "github.com/l-fraga2811/back-sable/internal/config"
    "github.com/l-fraga2811/back-sable/internal/models"
)

func main() {
    cfg := config.LoadConfig()

    // Executar auto-migration
    err := cfg.DB.AutoMigrate(
        &models.Item{},
    )

    if err != nil {
        log.Fatal("Erro na migração:", err)
    }

    log.Println("Migração concluída com sucesso!")
}