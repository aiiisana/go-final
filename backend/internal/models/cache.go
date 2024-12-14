package models

import (
	"time"

	"gorm.io/gorm"
)

type Cache struct {
	gorm.Model
	CacheKey       string    `json:"cache_key" db:"cache_key"`
	CacheValue     string    `json:"cache_value" db:"cache_value"`
	ExpirationTime time.Time `json:"expiration_time" db:"expiration_time"`
}
