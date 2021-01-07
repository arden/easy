package gf_plus

import (
	"github.com/arden/gf-plus/plus"
	gfplusRedis "github.com/arden/gf-plus/redis"
	"gorm.io/gorm"
)

// Redis returns an instance of redis client with specified configuration group name.
func Redis(name ...string) *gfplusRedis.Redis {
	return plus.GetRedis(name...)
}

// Gorm returns an instance of database ORM object with specified configuration group name.
func Gorm(name ...string) *gorm.DB {
	return plus.GetGorm(name...)
}
