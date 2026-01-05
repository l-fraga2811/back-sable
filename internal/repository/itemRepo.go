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