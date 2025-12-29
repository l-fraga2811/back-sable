# Migração para GORM mantendo Supabase

## Visão Geral

Este documento explica como migrar o projeto atual que usa Supabase via API REST para GORM mantendo a compatibilidade com o banco de dados Supabase PostgreSQL.

## Arquitetura Atual

- **Supabase Client**: Comunicação via API REST (PostgREST)
- **Models**: Structs Go sem tags de ORM
- **Handlers**: Lógica de negócio com chamadas HTTP para Supabase
- **Framework**: Fiber v3
- **Autenticação**: JWT tokens validados via middleware

## Arquitetura Futura

- **GORM**: ORM nativo para Go
- **Database Connection**: Conexão direta PostgreSQL
- **Models**: Structs Go com tags GORM
- **Repository Pattern**: Separação de concerns
- **Framework**: Mantido Fiber v3
- **Autenticação**: Mantida (JWT do Supabase)

## Passos da Migração

### 1. Preparação do Ambiente

#### 1.1 Dependências GORM (já instaladas)

```bash
# GORM já está no go.mod
gorm.io/gorm v1.31.1
gorm.io/driver/postgres v1.6.0
```

#### 1.2 Adicionar Variáveis de Ambiente

```env
# Variáveis existentes no projeto
PORT=3000
SUPABASE_URL=your_supabase_url
SUPABASE_KEY=your_supabase_anon_key
SUPABASE_JWKS_URL=your_supabase_jwks_url
SUPABASE_JWT_SECRET=your_supabase_jwt_secret
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_secret

# Novas variáveis para GORM
SUPABASE_DB_HOST=extracted_from_SUPABASE_URL
SUPABASE_DB_PORT=5432
SUPABASE_DB_USER=postgres
SUPABASE_DB_PASSWORD=your_db_password
SUPABASE_DB_NAME=postgres
```

### 2. Configuração do Database

#### 2.1 Criar Configuração GORM

```go
// config/database.go
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
```

#### 2.2 Atualizar Config.go

```go
// config/config.go
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
```

### 3. Atualizar Models com Tags GORM

#### 3.1 Model Item com Tags GORM

```go
// internal/models/item.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Item struct {
    ID          string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID      string         `gorm:"type:uuid;not null;index" json:"user_id"`
    Title       string         `gorm:"type:text;not null" json:"title"`
    Description string         `gorm:"type:text" json:"description"`
    Price       float64        `gorm:"type:decimal(12,2)" json:"price"`
    Completed   bool           `gorm:"default:false" json:"completed"`
    CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Item) TableName() string {
    return "items"
}

type CreateItemRequest struct {
    Title       string  `json:"title" validate:"required"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}

type UpdateItemRequest struct {
    Title       string  `json:"title,omitempty"`
    Description string  `json:"description,omitempty"`
    Price       float64 `json:"price,omitempty"`
    Completed   bool    `json:"completed,omitempty"`
}
```

### 4. Repository Pattern com GORM

#### 4.1 Interface Repository

```go
// internal/repository/item_repository.go
package repository

import "github.com/l-fraga2811/back-sable/internal/models"

type ItemRepository interface {
    Create(item *models.Item) error
    GetByID(id string) (*models.Item, error)
    GetAll(userID string) ([]models.Item, error)
    Update(item *models.Item) error
    Delete(id string) error
    GetByUserID(userID string) ([]models.Item, error)
}
```

#### 4.2 Implementação GORM

```go
// internal/repository/item_repository_gorm.go
package repository

import (
    "github.com/l-fraga2811/back-sable/internal/models"
    "gorm.io/gorm"
)

type itemRepositoryGORM struct {
    db *gorm.DB
}

func NewItemRepositoryGORM(db *gorm.DB) ItemRepository {
    return &itemRepositoryGORM{db: db}
}

func (r *itemRepositoryGORM) Create(item *models.Item) error {
    return r.db.Create(item).Error
}

func (r *itemRepositoryGORM) GetByID(id string) (*models.Item, error) {
    var item models.Item
    err := r.db.Where("id = ?", id).First(&item).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}

func (r *itemRepositoryGORM) GetAll(userID string) ([]models.Item, error) {
    var items []models.Item
    err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
    return items, err
}

func (r *itemRepositoryGORM) Update(item *models.Item) error {
    return r.db.Save(item).Error
}

func (r *itemRepositoryGORM) Delete(id string) error {
    return r.db.Delete(&models.Item{}, "id = ?", id).Error
}

func (r *itemRepositoryGORM) GetByUserID(userID string) ([]models.Item, error) {
    var items []models.Item
    err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
    return items, err
}
```

### 5. Atualizar Handlers para GORM

#### 5.1 ItemHandler com Repository Pattern

```go
// internal/handlers/item.go
package handlers

import (
    "github.com/gofiber/fiber/v3"
    "github.com/l-fraga2811/back-sable/internal/models"
    "github.com/l-fraga2811/back-sable/internal/repository"
)

type ItemHandler struct {
    itemRepo repository.ItemRepository
}

func NewItemHandler(itemRepo repository.ItemRepository) *ItemHandler {
    return &ItemHandler{
        itemRepo: itemRepo,
    }
}

func (h *ItemHandler) Create(c fiber.Ctx) error {
    userID := h.requireAuth(c)
    if userID == "" {
        return nil
    }

    var req models.CreateItemRequest
    if err := c.Bind().Body(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data: " + err.Error()})
    }

    item := &models.Item{
        UserID:      userID,
        Title:       req.Title,
        Description: req.Description,
        Price:       req.Price,
        Completed:   false,
    }

    if err := h.itemRepo.Create(item); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating item"})
    }

    return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *ItemHandler) GetAll(c fiber.Ctx) error {
    userID := h.requireAuth(c)
    if userID == "" {
        return nil
    }

    items, err := h.itemRepo.GetAll(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching items"})
    }

    return c.JSON(items)
}

func (h *ItemHandler) GetByID(c fiber.Ctx) error {
    userID := h.requireAuth(c)
    if userID == "" {
        return nil
    }

    itemID := c.Params("id")
    item, err := h.itemRepo.GetByID(itemID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
    }

    if item.UserID != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to access this item"})
    }

    return c.JSON(item)
}

func (h *ItemHandler) Update(c fiber.Ctx) error {
    userID := h.requireAuth(c)
    if userID == "" {
        return nil
    }

    itemID := c.Params("id")
    item, err := h.itemRepo.GetByID(itemID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
    }

    if item.UserID != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to update this item"})
    }

    var req models.UpdateItemRequest
    if err := c.Bind().Body(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data: " + err.Error()})
    }

    if req.Title != "" {
        item.Title = req.Title
    }
    if req.Description != "" {
        item.Description = req.Description
    }
    if req.Price != 0 {
        item.Price = req.Price
    }
    item.Completed = req.Completed

    if err := h.itemRepo.Update(item); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating item"})
    }

    return c.JSON(item)
}

func (h *ItemHandler) Delete(c fiber.Ctx) error {
    userID := h.requireAuth(c)
    if userID == "" {
        return nil
    }

    itemID := c.Params("id")
    item, err := h.itemRepo.GetByID(itemID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
    }

    if item.UserID != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to delete this item"})
    }

    if err := h.itemRepo.Delete(itemID); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting item"})
    }

    return c.JSON(fiber.Map{"message": "Item deleted successfully"})
}

func (h *ItemHandler) requireAuth(c fiber.Ctx) string {
    userID, ok := c.Locals("userID").(string)
    if !ok || userID == "" {
        c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
        return ""
    }
    return userID
}
```

### 6. Atualizar Main.go

#### 6.1 Nova Configuração do Main

```go
// cmd/api/main.go
package main

import (
    "log"

    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/cors"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/recover"
    "github.com/l-fraga2811/back-sable/internal/config"
    "github.com/l-fraga2811/back-sable/internal/handlers"
    "github.com/l-fraga2811/back-sable/internal/repository"
    "github.com/l-fraga2811/back-sable/internal/repository/supabase"
    "github.com/l-fraga2811/back-sable/internal/routes"
)

func main() {
    // Load Configuration
    cfg := config.LoadConfig()

    // Initialize Dependencies
    tokenValidator := supabase.NewTokenValidator(cfg)

    // Initialize GORM Repository
    itemRepo := repository.NewItemRepositoryGORM(cfg.DB)

    // Initialize Handlers
    itemHandler := handlers.NewItemHandler(itemRepo)

    // Initialize Fiber
    app := fiber.New(fiber.Config{
        AppName: "Sable Backend",
    })

    // Middleware
    app.Use(logger.New())
    app.Use(recover.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: []string{"*"},
        AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    }))

    // Setup Routes
    routes.SetupRoutes(app, tokenValidator, itemHandler)

    // Start Server
    log.Printf("Server starting on port %s", cfg.Port)
    log.Fatal(app.Listen(":" + cfg.Port))
}
```

### 7. Atualizar Rotas

#### 7.1 Nova Configuração de Rotas

```go
// internal/routes/routes.go
package routes

import (
    "github.com/gofiber/fiber/v3"
    "github.com/l-fraga2811/back-sable/internal/handlers"
    "github.com/l-fraga2811/back-sable/internal/repository/supabase"
)

func SetupRoutes(app *fiber.App, tokenValidator *supabase.TokenValidator, itemHandler *handlers.ItemHandler) {
    api := app.Group("/api")

    // Auth routes (mantidos)
    auth := api.Group("/auth")
    auth.Post("/signin", handlers.SignIn)
    auth.Post("/signup", handlers.SignUp)

    // Protected routes
    protected := api.Group("/")
    protected.Use(supabase.AuthMiddleware(tokenValidator))

    // Item routes (agora com GORM)
    items := protected.Group("/items")
    items.Get("/", itemHandler.GetAll)
    items.Post("/", itemHandler.Create)
    items.Get("/:id", itemHandler.GetByID)
    items.Put("/:id", itemHandler.Update)
    items.Delete("/:id", itemHandler.Delete)

    // Health check
    app.Get("/health", handlers.Health)
}
```

### 8. Migração dos Dados

#### 8.1 Script de Migração

```go
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
```

## Considerações Importantes

### 8.1 Compatibilidade Supabase

1. **Connection String**: Use as variáveis existentes do projeto
2. **Schema**: Mantenha o schema existente (não use `public.` prefix)
3. **UUIDs**: Configure GORM para usar UUIDs como no banco atual
4. **Timestamps**: Mantenha `timestamptz` para compatibilidade

### 8.2 Autenticação

- **Mantenha o JWT validation** do Supabase Auth
- **Use o user_id** do token para filtrar dados
- **Não altere o middleware** de autenticação existente

### 8.3 Performance

- **Connection Pool**: GORM gerencia automaticamente
- **Queries**: Use `preload` para relacionamentos
- **Indexes**: Mantenha os índices existentes

## Benefícios da Migração

1. **Performance**: Conexão direta vs HTTP
2. **Type Safety**: Compile-time query validation
3. **Features**: Transactions, migrations, relationships
4. **Debugging**: SQL logging integrado
5. **Testing**: Mock repositories facilmente

## Riscos e Mitigações

### Riscos:

- **Downtime** durante migração
- **Data loss** se migration falhar
- **Performance regressions** temporárias

### Mitigações:

- **Backup** completo antes de iniciar
- **Teste** em ambiente staging
- **Rollback plan** preparado
- **Monitoramento** pós-migração

## Checklist Final

- [ ] Backup do banco de dados
- [ ] Testar em ambiente staging
- [ ] Configurar variáveis de ambiente GORM
- [ ] Criar config/database.go
- [ ] Atualizar config/config.go
- [ ] Adicionar tags GORM aos models
- [ ] Implementar repository pattern
- [ ] Atualizar handlers para usar repositories
- [ ] Atualizar main.go
- [ ] Atualizar routes
- [ ] Executar migração
- [ ] Testar todos os endpoints
- [ ] Performance testing
- [ ] Deploy para produção
- [ ] Monitorar por 24h

## Comandos Úteis

```bash
# Testar conexão
go run cmd/test-db/main.go

# Executar migração
go run cmd/migrate/main.go

# Build para produção
go build -o api cmd/api/main.go

# Executar servidor
./api
```

## Recursos

- [GORM Documentation](https://gorm.io/docs/)
- [Supabase PostgreSQL](https://supabase.com/docs/guides/database/connecting-to-postgres)
- [Fiber Documentation](https://docs.gofiber.io/)
