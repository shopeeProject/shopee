package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	models "github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/storage"
	util "github.com/shopeeProject/shopee/util"
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

func getUserDB() *gorm.DB {

	db := getStorageConfig()
	err := models.MigrateUser(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}
	return db

}

func getSellerDB() *gorm.DB {

	db := getStorageConfig()
	err := models.MigrateSeller(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}
	return db

}

func main() {
	server := NewAPIServer(":3000") // runs on 3000
	shopeeDB := util.ShopeeDatabase{
		UserDB:   getUserDB(),
		SellerDB: getSellerDB(),
	}
	server.Run(&shopeeDB)

	// call run
	fmt.Println("Hi Buddy!!, Server is running")
}
