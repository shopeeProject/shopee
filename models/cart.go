package models

import "gorm.io/gorm"

type Cart struct {
	CartID uint `gorm:"primaryKey;autoIncrement" json:"cartid"`
	UID    int  ` json:"uid"`
	PID    int  `json:"pid"`
	Count  int  `json:"count"`
}

func MigrateCart(db *gorm.DB) error {
	err := db.AutoMigrate(&Cart{})
	return err
}
