package models

import (
	"gorm.io/gorm"
)

type Seller struct {
	SID          uint   `gorm:"primaryKey;autoIncrement" json:"sid"`
	Name         string `json:"name"`
	EmailAddress string `json:"email"`
	Rating       uint   `json:"rating"`
	Password     string `json:"password"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Status       string `json:"status"`
}

func MigrateSeller(db *gorm.DB) error {
	err := db.AutoMigrate(&Seller{})
	return err
}
