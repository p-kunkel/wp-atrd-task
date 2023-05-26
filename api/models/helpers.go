package models

import "gorm.io/gorm"

func mustHaveRecord(DB *gorm.DB) error {
	if DB.Error != nil {
		return DB.Error
	}
	if DB.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
