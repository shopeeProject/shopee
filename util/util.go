package util

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}
