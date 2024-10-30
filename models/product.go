package models

import (
	"gorm.io/gorm"
)

type Product struct {
	PID          int     `gorm:"primary key;autoIncrement" json:"pid"`
	Name         string  `json:"name"`
	Price        int     `json:"price"`
	Availability bool    `json:"availability"`
	Rating       float32 `json:"rating" gorm:"default:0"`
	CategoryID   int     `json:"category"`
	Count        int     `json:"count"`
	Description  string  `json:"description"`
	SID          string  `json:"sid"`
	Image        string  `json:"image"`
}

/*

JWT authentication
Product creation and methods
Cart creation and methods
Payment
Order


*/

func MigrateProduct(db *gorm.DB) error {
	err := db.AutoMigrate(&Product{})
	return err
}
