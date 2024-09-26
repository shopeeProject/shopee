package models

import "gorm.io/gorm"

type Category struct {
	CategoryId uint   `gorm:"primaryKey;autoIncrement" json:"cat_id"`
	Id         uint   ` json:"id"`
	Name       string `json:"name"`
}

func MigrateCategory(db *gorm.DB) error {
	err := db.AutoMigrate(&Category{})
	return err
}
