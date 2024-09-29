package models

import (
	"gorm.io/gorm"
)

type Seller struct {
	SID          int    `gorm:"primaryKey;autoIncrement" json:"sid"`
	Name         string `json:"name"`
	EmailAddress string `gorm:"unique;required" json:"email"`
	Rating       int    `json:"rating"`
	Password     string `json:"password"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Status       string `json:"status"`
}

func MigrateSeller(db *gorm.DB) error {
	err := db.AutoMigrate(&Seller{})
	return err
}
