package data

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandInt returns a random int in [0, n).
func RandInt(n int) int {
	if n <= 0 {
		return 0
	}
	return rand.Intn(n)
}
