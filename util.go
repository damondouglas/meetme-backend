package meetme

import (
	"math/rand"
	"time"
)

var (
	// source = rand.Seed(time.Now().UTC().UnixNano())
	// r = rand.New(source)
	r = NewRandom()
)

func contains(arrayOfStrings []string, s string) bool {
	for _, a := range arrayOfStrings {
		if a == s {
			return true
		}
	}
	return false
}

func randbetween(min int, max int) int {
	return r.Intn(max-min) + min
}

// NewRandom creates new [rand.Rand] with UTC.UnixNano seed
func NewRandom() *rand.Rand {
	t := time.Now().UTC().UnixNano()
	source := rand.NewSource(t)
	return rand.New(source)
}
