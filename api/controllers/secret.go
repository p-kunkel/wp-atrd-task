package controllers

import (
	"log"
	"server/config"
	"server/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func AddSecret(c *gin.Context) {
	var (
		s   models.Secret
		err error
	)

	if err = c.MustBindWith(&s, binding.FormPost); err != nil {
		handleError(c, newErrResp(405, err))
		return
	}

	if err = s.CheckAndFill(); err != nil {
		handleError(c, newErrResp(405, err))
		return
	}

	if err = s.Insert(config.DB); err != nil {
		handleError(c, newErrResp(405, err))
		return
	}

	c.JSON(200, s)
}

func GetSecret(c *gin.Context) {
	var (
		s   models.Secret
		err error
	)

	if s.Hash = c.Param("hash"); s.Hash == "" {
		handleError(c, newErrResp(404, err))
		return
	}

	log.Println(s.Hash)
	if err = s.GetIfCanBeTaken(config.DB); err != nil {
		handleError(c, newErrResp(404, err))
		return
	}

	c.JSON(200, s)
}
