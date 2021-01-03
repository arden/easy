package example

import (
	"context"
	"github.com/arden/gf-plus/plus"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	redisClient := plus.Redis().GetClient()
	redisClient.Set(context.Background(), "hdcj:test", "arden", time.Hour)
}
