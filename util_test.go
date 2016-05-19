package meetme

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandom(t *testing.T) {
	count := 10
	min := 10
	max := 100
	for i := 0; i < count; i++ {
		v := randbetween(min, max)
		assert.True(t, v >= min, "Random number should be greater than or equal to %s", min)
		assert.True(t, v <= max, "Random number should be less than or equal to %s", max)
	}
}
