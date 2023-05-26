package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	if !assert.NoError(t, LoadEnv("../.env")) {
		t.FailNow()
	}

	if !assert.NoError(t, ConnectDB(DB, GetDBAddress())) {
		t.FailNow()
	}
}

func TestLoadEnv(t *testing.T) {
	if !assert.NoError(t, LoadEnv("../.env")) {
		t.FailNow()
	}
}
