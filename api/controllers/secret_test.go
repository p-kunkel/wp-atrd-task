package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"server/config"
	"server/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	failIfErr(models.AutoMigrateDB(config.DB))
}

func TestAddSecret(t *testing.T) {
	s := models.Secret{}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/v1/secret", AddSecret)

	body := url.Values{
		"secret":           []string{"test"},
		"expireAfterViews": []string{"0"},
		"expireAfter":      []string{"-1"},
	}

	f := func(b url.Values, expectedStatusCode int, respBody *models.Secret) {
		req, err := http.NewRequest(http.MethodPost, "/v1/secret", bytes.NewBufferString(b.Encode()))
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if !assert.Equal(t, expectedStatusCode, w.Code) {
			t.FailNow()
		}

		if expectedStatusCode == http.StatusOK {
			if !assert.NoError(t, json.NewDecoder(w.Body).Decode(&respBody)) {
				t.FailNow()
			}
		}
	}

	f(body, http.StatusMethodNotAllowed, &s)
	body["expireAfterViews"] = []string{"1"}
	f(body, http.StatusMethodNotAllowed, &s)
	body["expireAfter"] = []string{"0"}
	f(body, http.StatusOK, &s)

	if !assert.NotEmpty(t, s.Hash) {
		t.FailNow()
	}
}

func TestGetSecret(t *testing.T) {
	h := uuid.NewString()
	s := models.Secret{
		Hash:           h,
		ExpiresAt:      0,
		RemainingViews: 1,
		SecretText:     "test",
	}

	if !assert.NoError(t, s.Insert(config.DB)) {
		t.FailNow()
	}
	s = models.Secret{}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/v1/secret/:hash", GetSecret)

	f := func(hash string, expectedStatusCode int, respBody *models.Secret) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/secret/%s", hash), nil)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if !assert.Equal(t, expectedStatusCode, w.Code) {
			t.FailNow()
		}

		if expectedStatusCode == http.StatusOK {
			if !assert.NoError(t, json.NewDecoder(w.Body).Decode(&respBody)) {
				t.FailNow()
			}
		}
	}

	f("", http.StatusNotFound, &s)
	f(h, http.StatusOK, &s)

	if !assert.Equal(t, h, s.Hash) {
		t.FailNow()
	}

	f(h, http.StatusNotFound, &s)
}
