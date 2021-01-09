package example

import (
	"context"
	"github.com/arden/easy"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	redisClient := easy.Redis()
	redisClient.Set(context.Background(), "hdcj:test", "11arden", time.Hour)
}
