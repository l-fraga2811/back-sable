package repository

import "github.com/l-fraga2811/back-sable/internal/models"

type ProfileRepository interface {
	GetByID(id string) (*models.Profile, error)
}
