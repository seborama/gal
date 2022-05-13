package gal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	xpn := `-3 + 4`
	val := Eval(xpn)
	assert.Equal(t, NewNumber(1), val)
}
