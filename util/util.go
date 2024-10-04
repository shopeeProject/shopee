package util

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

type Response struct {
	Success bool
	Message string
}

type DataResponse struct {
	Success bool
	Message string
	Data    map[string]string
}
