package util

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
)

const IntMax = int(^uint(0) >> 1)
const IntMin = ^IntMax

// RandomMinMax 生成指定范围内的随机数字
func RandomMinMax(min, max int64) int64 {
	if min > max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

// RandomInt64 生成指定位数的随机数字
func RandomInt64(b ...int) int64 {
	l := 5
	if len(b) > 0 && b[0] != 0 {
		l = b[0]
	}
	w := math.Pow(10, float64(l))
	n := int64(0)
	if w > math.MaxFloat64 {
		n = int64(math.MaxInt64)
	} else {
		n = int64(w)
	}
	result, _ := crand.Int(crand.Reader, big.NewInt(n))
	if l > len(result.String()) {
		return RandomInt64(l)
	}
	return result.Int64()
}
