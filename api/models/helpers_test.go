package models

import (
	"log"
	"server/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	failIfErr := func(err error) {
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	config.SetTimeUTC()
	failIfErr(config.LoadEnv("../.env"))

	failIfErr(config.ConnectDB(config.DB, config.GetDBAddress()))
	config.DB.Config.Logger = logger.Default.LogMode(logger.Error)
	failIfErr(AutoMigrateDB(config.DB))
}
func TestMustHaveRecord(t *testing.T) {
	err := mustHaveRecord(config.DB.Where("hash = ''").Find(&Secret{}))
	if !assert.EqualError(t, err, gorm.ErrRecordNotFound.Error()) {
		t.FailNow()
	}
}
