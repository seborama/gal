package gal_test

import (
	"testing"

	"github.com/seborama/gal"
	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	xpn := `-3 + 4`
	val := gal.Parse(xpn).Eval()
	assert.Equal(t, gal.NewNumber(1), val)
}
