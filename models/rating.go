package models

import "gorm.io/gorm"

type Rating struct {
	RID         int    `gorm:"primary key;autoIncrement" json:"rid"`
	UID         int    `json:"uid"`
	PID         int    `json:"pid"`
	Rating      string `json:"rating"`
	RatingValue int    `json:"ratingValue"`
	Description string `json:"description"`
}

func MigrateRating(db *gorm.DB) error {
	err := db.AutoMigrate(&Rating{})
	return err
}
