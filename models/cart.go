package models

import "gorm.io/gorm"

type Cart struct {
	CartID uint `gorm:"primaryKey;autoIncrement" json:"cartid"`
	UID    uint ` json:"uid"`
	PID    uint `json:"pid"`
	Count  uint `json:"count"`
}

func MigrateCart(db *gorm.DB) error {
	err := db.AutoMigrate(&Cart{})
	return err
}
