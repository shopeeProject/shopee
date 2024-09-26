package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	PayID         int    `gorm:"primary key;autoIncrement" json:"rid"`
	UID           int    `json:"uid"`
	Amount        int    `json:"amount"`
	Timestamp     string `json:"timestamp"`
	PaymentStatus int    `json:"paymentStatus"`
	Description   string `json:"description"`
}

func MigratePayment(db *gorm.DB) error {
	err := db.AutoMigrate(&Payment{})
	return err
}
