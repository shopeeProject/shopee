package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	models "github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/storage"
	"gorm.io/gorm"
)

func getStorageConfig() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	return db
}

func getUserDB(db *gorm.DB) {

	err := models.MigrateUser(db)
	if err != nil {
		log.Fatal("could not migrate User db")
	}

}

func getSellerDB(db *gorm.DB) {
	err := models.MigrateSeller(db)
	if err != nil {
		log.Fatal("could not migrate Seller db")
	}

}

func getCartDB(db *gorm.DB) *gorm.DB {
	err := models.MigrateCart(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}
	return db

}

func getProductDB(db *gorm.DB) {
	err := models.MigrateProduct(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

}

func getOrderDB(db *gorm.DB) {
	err := models.MigrateOrder(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

}

func getCategoryDB(db *gorm.DB) {
	err := models.MigrateCategory(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

}

func migrateAllDB(db *gorm.DB) {
	getUserDB(db)
	getSellerDB(db)
	getCartDB(db)
	getOrderDB(db)
	getCategoryDB(db)
	getProductDB(db)

}

func main() {
	server := NewAPIServer(":5000") // runs on 5000

	db := getStorageConfig()
	migrateAllDB(db)
	server.Run(db)

	fmt.Println("Hi Buddy!!, Server is running")
}
