package main

import (
	"context"
	"github.com/arden/gf-plus/plus"
	"time"
)

func main()  {
	redisClient := plus.Redis().GetClient()
	redisClient.Set(context.Background(), "hdcj:test", "arden", time.Hour)
}
