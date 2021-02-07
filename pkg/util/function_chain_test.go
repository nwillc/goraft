package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSingleChainForward(t *testing.T) {
	counter := 2
	fc := FunctionChain{}
	fc.Add(func() {
		if counter == 2 {
			counter--
		}
	}).Add(func() {
		counter--
	})
	assert.Equal(t, 2, counter)
	fc.InvokeForward()
	assert.Equal(t, 0, counter)
}

func TestSingleChainRevers(t *testing.T) {
	counter := 2
	fc := FunctionChain{}
	fc.Add(func() {
		if counter == 1 {
			counter--
		}
	}).Add(func() {
		counter--
	})
	assert.Equal(t, 2, counter)
	fc.InvokeReverse()
	assert.Equal(t, 0, counter)
}
