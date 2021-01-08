package example

import (
	"context"
	"github.com/arden/gf-plus"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	redisClient := gf_plus.Redis()
	redisClient.Set(context.Background(), "hdcj:test", "11arden", time.Hour)
}
