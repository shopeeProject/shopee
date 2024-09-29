package models

import (
	"gorm.io/gorm"
)

type User struct {
	UId           int    `gorm:"primary key;autoIncrement" json:"id"`
	Name          string `json:"name"`
	PhoneNumber   string `json:"phoneNumber"`
	EmailAddress  string `gorm:"unique;required" json:"emailAddress"`
	AccountStatus string `json:"accountStatus"`
	Address       string `json:"address"`
	Password      string `json:"password"`
}

func MigrateUser(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
