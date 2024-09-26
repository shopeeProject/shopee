package models

import "gorm.io/gorm"

type Product struct {
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
