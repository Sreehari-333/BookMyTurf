package db

import (
	"BookMyTurf/models"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connecting to the database

func ConnectDB() {

	var err error

	err = godotenv.Load()

	if err != nil {
		log.Fatal("failed to load env")
	}

	dsn := os.Getenv("DB_DSN")

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database connection failed", err)
	}

	log.Println("Database connected successfully")

	// Migrating tables

	err = DB.AutoMigrate(&models.User{}, &models.Turf{}, &models.Booking{}, &models.RefreshToken{}, &models.BlockedSlot{}, &models.Payment{})

	if err != nil {
		log.Fatal("error with creating table", err)
	}

	// connection pooling

	sqlDB, err := DB.DB()

	if err != nil {
		log.Fatal("failed to get sql.db from gorm.db")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
}
