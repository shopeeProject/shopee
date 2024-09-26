package models

import (
	"gorm.io/gorm"
)

type Order struct {
	OID           int      `gorm:"primaryKey;autoIncrement" json:"oid"`
	UID           int      `json:"uid"`
	Price         int      `json:"price"`
	OrderStatus   string   `json:"order_status"`
	PaymentID     int      `json:"payment_id"`
	ProductsLists []int    `gorm:"type:bigint[]" json:"products_lists"`
	PaymentStatus string   `json:"payment_status"`
	Address       string   `json:"address"`
	StagesList    []string `gorm:"type:text[]"json:"stages_list"`
}

func MigrateOrder(db *gorm.DB) error {
	err := db.AutoMigrate(&Order{})
	return err
}
