package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	models "github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/storage"
	util "github.com/shopeeProject/shopee/util"
)

func main() {
	server := NewAPIServer(":3000") // runs on 3000

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
	err = models.MigrateUser(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := util.Repository{
		DB: db,
	}
	server.Run(&r)

	// call run
	fmt.Println("Hi Buddy!!, Server is running")
}
