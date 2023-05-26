package main

import (
	"log"
	"server/config"
	"server/mappings"
	"server/models"
)

func main() {
	failIfErr := func(err error) {
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	config.SetTimeUTC()
	failIfErr(config.LoadEnv())

	failIfErr(config.ConnectDB(config.DB, config.GetDBAddress()))
	failIfErr(models.AutoMigrateDB(config.DB))

	failIfErr(mappings.RunServer())
}
