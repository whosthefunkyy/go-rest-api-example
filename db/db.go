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


func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ConnectGorm() {
	
	host := getEnv("RDS_HOSTNAME", getEnv("DB_HOST", "localhost"))
	port := getEnv("RDS_PORT", getEnv("DB_PORT", "5432"))
	user := getEnv("RDS_USERNAME", getEnv("DB_USER", "artem"))
	pass := getEnv("RDS_PASSWORD", getEnv("DB_PASSWORD", "password"))
	name := getEnv("RDS_DB_NAME", getEnv("DB_NAME", "ebdb"))

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name,
	)

	fmt.Printf("Connecting to host: %s, database: %s\n", host, name)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
}
