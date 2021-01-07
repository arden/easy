package gf_plus

import (
	"github.com/arden/gf-plus/plus"
	gfplusRedis "github.com/arden/gf-plus/redis"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Redis returns an instance of redis client with specified configuration group name.
func Redis(name ...string) *gfplusRedis.Redis {
	return plus.GetRedis(name...)
}

func RedisClient(name ...string) *redis.Client {
	return plus.GetRedis(name...).GetClient()
}

// Gorm returns an instance of database ORM object with specified configuration group name.
func Gorm(name ...string) *gorm.DB {
	return plus.GetGorm(name...)
}
