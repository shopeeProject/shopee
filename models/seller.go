package models

import (
	"gorm.io/gorm"
)

type seller struct {
	SID    				int    `gorm:"primary key;autoIncrement" json:"id"`
	name   				*string
	emailAddress 		*string
	rating 				*int
	description  		*string 
	image  				*string
	status				*string

}

func MigrateSeller(db *gorm.DB) error {
	err := db.AutoMigrate(&seller{})
	return err
}