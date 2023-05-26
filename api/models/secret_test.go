package models

import (
	"log"
	"server/config"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ErrRollBack string = "sql: transaction has already been committed or rolled back"

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

func TestSecretCheckAndFill(t *testing.T) {
	s := Secret{
		ExpireAfter:    1,
		RemainingViews: 1,
	}

	if !assert.NoError(t, s.CheckAndFill()) {
		t.FailNow()
	}

	if !assert.NotEmpty(t, s.Hash) {
		t.FailNow()
	}
}

func TestSecretValid(t *testing.T) {
	s := Secret{
		Hash:           uuid.NewString(),
		ExpireAfter:    0,
		RemainingViews: 1,
	}

	if !assert.NoError(t, s.Valid()) {
		t.FailNow()
	}

	s.ExpireAfter = -1
	if err := s.Valid(); err == nil {
		assert.FailNow(t, "expected error, got nill")
	}

	s.ExpireAfter = 2
	s.RemainingViews = 0
	if err := s.Valid(); err == nil {
		assert.FailNow(t, "expected error, got nill")
	}
}

func TestSecretInsert(t *testing.T) {
	expected := Secret{
		Hash:           uuid.NewString(),
		ExpiresAt:      TimePointer(time.Now().Add(3 * time.Second)),
		RemainingViews: 1,
		SecretText:     "test",
	}
	actual := Secret{Hash: expected.Hash}

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		if err := expected.Insert(tx); err != nil {
			return err
		}

		if err := actual.GetIfCanBeTaken(tx); err != nil {
			return err
		}

		return tx.Rollback().Error
	})

	if !assert.EqualError(t, err, ErrRollBack) {
		t.FailNow()
	}

	if !assert.Equal(t, expected.SecretText, actual.SecretText) ||
		!assert.Equal(t, expected.RemainingViews-1, actual.RemainingViews) {
		t.FailNow()
	}
}

func TestSecretGetIfCanBeTaken(t *testing.T) {
	expected := Secret{
		Hash:           uuid.NewString(),
		ExpiresAt:      TimePointer(time.Now().Add(3 * time.Second)),
		RemainingViews: 1,
		SecretText:     "test",
	}
	actual := Secret{Hash: expected.Hash}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		expected.Insert(tx)

		if err := actual.GetIfCanBeTaken(tx); err != nil {
			return err
		}

		if err := actual.GetIfCanBeTaken(tx); err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return err
		}

		expected.Hash = uuid.NewString()
		expected.ExpiresAt = TimePointer(time.Now().Add(2 * time.Second))
		expected.RemainingViews = 2
		expected.Insert(tx)

		actual.Hash = expected.Hash
		if err := actual.GetIfCanBeTaken(tx); err != nil {
			return err
		}

		time.Sleep(3 * time.Second)
		if err := actual.GetIfCanBeTaken(tx); err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return err
		}

		return tx.Rollback().Error
	})

	if !assert.EqualError(t, err, ErrRollBack) {
		t.FailNow()
	}

	if !assert.Equal(t, expected.SecretText, actual.SecretText) ||
		!assert.Equal(t, expected.RemainingViews-1, actual.RemainingViews) {
		t.FailNow()
	}
}

func TestSecretFindByHash(t *testing.T) {
	expected := Secret{
		Hash:           uuid.NewString(),
		ExpiresAt:      TimePointer(time.Now().Add(time.Minute)),
		RemainingViews: 1,
		SecretText:     "test",
	}
	actual := Secret{Hash: expected.Hash}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		expected.Insert(tx)

		if err := actual.FindByHash(tx); err != nil {
			return err
		}

		return tx.Rollback().Error
	})

	if !assert.EqualError(t, err, ErrRollBack) {
		t.FailNow()
	}

	if !assert.Equal(t, expected.SecretText, actual.SecretText) {
		t.FailNow()
	}

}

func TestSecretFind(t *testing.T) {
	expected := Secret{
		Hash:           uuid.NewString(),
		ExpiresAt:      TimePointer(time.Now().Add(time.Minute)),
		RemainingViews: 1,
		SecretText:     "test",
	}
	actual := Secret{}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		scope := func(db *gorm.DB) *gorm.DB { return db.Where("hash = ?", expected.Hash) }
		expected.Insert(tx)

		if err := actual.Find(tx, scope); err != nil {
			return err
		}

		return tx.Rollback().Error
	})

	if !assert.EqualError(t, err, ErrRollBack) {
		t.FailNow()
	}

	if !assert.Equal(t, expected.Hash, actual.Hash) ||
		!assert.Equal(t, expected.SecretText, actual.SecretText) {
		t.FailNow()
	}
}

func TestSecretTableName(t *testing.T) {
	s := Secret{}
	if !assert.NotEmpty(t, s.TableName(), "empty name of secret table") {
		t.FailNow()
	}
}
