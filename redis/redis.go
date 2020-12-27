package redis

import (
	"github.com/go-redis/redis/v8"
)

const (
	frameCoreComponentNameRedis = "gfplus.core.component.redis"
	configNodeNameRedis         = "redis"
)

func Redis(name ...string) *redis.Client {
	return nil
}