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