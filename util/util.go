package util

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

type ShopeeDatabase struct {
	UserDB   *gorm.DB
	SellerDB *gorm.DB
}
