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
