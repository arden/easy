package plus

import (
	gfplusRedis "github.com/arden/gf-plus/redis"
	"gorm.io/gorm"
)

// Redis returns an instance of redis client with specified configuration group name.
func Redis(name ...string) *gfplusRedis.Redis {
	return GetRedis(name...)
}

// Gorm returns an instance of database ORM object with specified configuration group name.
func Gorm(name ...string) *gorm.DB {
	return GetGorm(name...)
}
