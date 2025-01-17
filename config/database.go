package config

import (
	"arfdev/chat/internal/domain/entities"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	SSLMode    string
}

func NewDB() *gorm.DB {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	config := dbConfig{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASS"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
		SSLMode:    "disable",
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
		config.SSLMode,
		loc,
	)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Println("[db] - Postgresql failed to connect. Reason:", err)
	} else {
		log.Println("[db] - Connected to Postgresql")
	}

	// Uncomment this after defining models
	tables := []interface{}{
		&entities.Mst_otp{},
		&entities.Mst_users{},
		&entities.Mst_users_detail{},
	}
	migrate(DB, tables)

	return DB
}

func migrate(db *gorm.DB, models []interface{}) {
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatal("[db] - Migration failed.\nReason:", err)
	} else {
		log.Println("[db] - Migration success.")
	}
}
