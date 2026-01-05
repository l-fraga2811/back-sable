package repository

import (
	"github.com/l-fraga2811/back-sable/internal/models"
	"gorm.io/gorm"
)

type profileRepositoryGorm struct {
	db *gorm.DB
}

func NewProfileRepositoryGorm(db *gorm.DB) ProfileRepository {
	return &profileRepositoryGorm{db: db}
}

func (r *profileRepositoryGorm) GetByID(id string) (*models.Profile, error) {
	var profile models.Profile
	if err := r.db.Where("id = ?", id).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}
