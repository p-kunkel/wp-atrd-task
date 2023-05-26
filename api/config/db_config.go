package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB = &gorm.DB{}

func ConnectDB(db *gorm.DB, address string) error {
	dbTemp, err := gorm.Open(postgres.Open(address), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return err
	}

	*db = *dbTemp
	return nil
}

func GetDBAddress() string {
	var (
		host     = os.Getenv("DB_ADDRESS")
		user     = os.Getenv("DB_LOGIN")
		password = os.Getenv("DB_PASSWORD")
		dbName   = os.Getenv("DB_NAME")
		port     = os.Getenv("DB_PORT")
	)

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbName, port)
}

func LoadEnv(filenames ...string) error {
	return godotenv.Load(filenames...)
}
