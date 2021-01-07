package example

import (
	"context"
	"github.com/arden/gf-plus"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	redisClient := gf_plus.Redis().GetClient()
	redisClient.Set(context.Background(), "hdcj:test", "arden", time.Hour)
}
