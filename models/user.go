package models

import (
	"gorm.io/gorm"
)

type user struct {
	UID    			uint    `gorm:"primary key;autoIncrement" json:"id"`
	name   			*string
	phoneNumber 	*string
	emailAddress 	*string
	accountStatus 	*string
	address  		*string 
}

func MigrateUser(db *gorm.DB) error {
	err := db.AutoMigrate(&user{})
	return err
}