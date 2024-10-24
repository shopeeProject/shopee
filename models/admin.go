package models

import (
	"gorm.io/gorm"
)

type Admin struct {
	UId          int    `gorm:"primary key;autoIncrement" json:"id"`
	Name         string `json:"name"`
	PhoneNumber  string `json:"phoneNumber"`
	EmailAddress string `gorm:"unique;required" json:"emailAddress"`
	Password     string `json:"password"`
}

func MigrateAdmin(db *gorm.DB) error {
	err := db.AutoMigrate(&Admin{})
	return err
}
