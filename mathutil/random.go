package mathutil

import (
	"math/rand"
	"time"
)

// RandomInt return an random int at the min ~ max
// Usage:
// 	RandomInt(10, 99)
// 	RandomInt(100, 999)
// 	RandomInt(1000, 9999)
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func RandomInt64(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Int63n(max-min)
}
