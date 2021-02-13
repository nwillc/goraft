package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLooper(t *testing.T) {
	count := 0
	stopper := RepeatUntilStopped(10*time.Millisecond, func() {
		count++
	})
	assert.NotNil(t, stopper)
	time.Sleep(100 * time.Millisecond)
	assert.LessOrEqual(t, 8, count)
	stopper <- true
	c1 := count
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, c1, count)
}
