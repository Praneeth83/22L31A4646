package config

import (
	"Backend-Test-Submission/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=praneeth dbname=urlshortner port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	if err := DB.AutoMigrate(&models.ShortURL{}, &models.ClickInfo{}); err != nil {
		log.Fatal("Failed to AutoMigrate")
	}
	fmt.Println("Connected to postgres server using GORM")
}
