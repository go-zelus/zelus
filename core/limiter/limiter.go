package limiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type ILimiter interface {
	Key(ctx *gin.Context) string
	GetBucket(key string) (*ratelimit.Bucket, bool)
	AddBuckets(rules ...BucketRule) ILimiter
}

type BucketRule struct {
	Key          string
	FillInterval time.Duration
	Capacity     int64
	Quantum      int64
}

type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}
