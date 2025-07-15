package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/whosthefunkyy/go-rest-api-example/models"
)

var DB *gorm.DB

func ConnectGorm() {

	  dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
}
func AutoMigrate() {
	err := DB.AutoMigrate(&models.User{}) 
	if err != nil {
		log.Fatalf("migration failed: %s", err)
	}
}