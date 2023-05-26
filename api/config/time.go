package config

import (
	"time"
)

func SetTimeUTC() {
	time.Local = time.UTC
}
