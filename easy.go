package easy

import (
	"github.com/arden/easy/database"
	easyRedis "github.com/arden/easy/database/redis"
	"gorm.io/gorm"
)

// Redis returns an instance of redis client with specified configuration group name.
func Redis(name ...string) *easyRedis.Redis {
	return database.GetRedis(name...)
}

// Gorm returns an instance of database ORM object with specified configuration group name.
func Gorm(name ...string) *gorm.DB {
	return database.GetGorm(name...)
}
