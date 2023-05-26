package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Secret struct {
	Hash           string    `json:"hash" form:"-" gorm:"type:varchar(36);notnull;unique;default:null"`
	CreatedAt      time.Time `json:"createdAt" form:"-" gorm:"type:timestamp;notnull"`
	ExpiresAt      int32     `json:"expiresAt" form:"expireAfter" gorm:"type:int8;notnull"`
	RemainingViews int32     `json:"remainingViews" form:"expireAfterViews" gorm:"type:int4;notnull;default:null"`
	SecretText     string    `json:"secretText" form:"secret" gorm:"type:varchar"`
}

func (s *Secret) CheckAndFill() error {
	if err := s.Valid(); err != nil {
		return err
	}

	s.Hash = uuid.NewString()
	if s.ExpiresAt > 0 {
		s.ExpiresAt = int32(time.Now().Add(time.Duration(s.ExpiresAt) * time.Minute).Unix())
	}

	return nil
}

func (s *Secret) Valid() error {
	if s.ExpiresAt < 0 {
		return errors.New("invalid expiresAfter value")
	}

	if s.RemainingViews <= 0 {
		return errors.New("invalid expireAfterViews value")
	}

	return nil
}

func (s *Secret) Insert(DB *gorm.DB) error {
	return DB.Create(&s).Error
}

func (s *Secret) GetIfCanBeTaken(DB *gorm.DB, scope ...func(*gorm.DB) *gorm.DB) error {
	return mustHaveRecord(
		DB.Model(&s).
			Where("hash = ? AND remaining_views > 0 AND (expires_at >= ? OR expires_at = 0)", s.Hash, time.Now().Unix()).
			Where("remaining_views > 0").
			Where("expires_at >= ? OR expires_at = 0", time.Now().Unix()).
			Clauses(clause.Returning{
				Columns: []clause.Column{
					{
						Name: "*",
						Raw:  true,
					},
				},
			}).UpdateColumn("remaining_views", gorm.Expr("remaining_views - 1")))
}

func (s *Secret) FindByHash(DB *gorm.DB, scope ...func(*gorm.DB) *gorm.DB) error {
	scope = append(scope, func(db *gorm.DB) *gorm.DB { return db.Where("hash = ?", s.Hash) })
	return s.Find(DB, scope...)
}

func (s *Secret) Find(DB *gorm.DB, scope ...func(*gorm.DB) *gorm.DB) error {
	return DB.Scopes(scope...).Find(&s).Error
}

func (*Secret) TableName() string {
	return "public.secret"
}
