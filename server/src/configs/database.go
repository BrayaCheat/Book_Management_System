package configs

import (
	"fmt"
	"log"
	"os"
	"server/src/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func AutoMigrate() error {
	if DB != nil {
		DB.AutoMigrate(
			&models.Author{},
			&models.Book{},
			&models.AuthorAddress{},
			&models.User{},		)
	} else {
		log.Fatal("Fail to init db")
	}
	return nil
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func InitDatabase() error {
	LoadEnv()
	usr := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PASSWORD")
	dbn := os.Getenv("DB_NAME")

	if usr == "" || pwd == "" || dbn == "" {
		return fmt.Errorf("one or more environment variables are missing")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(host.docker.internal:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", usr, pwd, dbn)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	err = AutoMigrate()
	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return nil
}
