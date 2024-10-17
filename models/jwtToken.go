package models

import "gorm.io/gorm"

type Token struct {
	RefreshToken string `gorm:"primaryKey" json:"refreshToken"`
	Email        string ` json:"emailAddress"`
}

func MigrateToken(db *gorm.DB) error {
	err := db.AutoMigrate(&Token{})
	return err
}