package utils

import (
	"math/rand"
	"time"
)

func Rand(n int64) int64 {
	value := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(n)
	return value
}
