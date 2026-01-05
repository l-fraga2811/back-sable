package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username   string    `gorm:"type:text" json:"username"`
	Phone      string    `gorm:"type:text" json:"phone"`
	ProfileUrl string    `gorm:"type:text;column:profile_url" json:"profileUrl"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Profile) TableName() string {
	return "profiles"
}
