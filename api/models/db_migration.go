package models

import "gorm.io/gorm"

func AutoMigrateDB(DB *gorm.DB) error {
	return DB.AutoMigrate(&Secret{})
}
