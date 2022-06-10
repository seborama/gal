package gal_test

import (
	"math"
	"testing"

	"github.com/seborama/gal"
	"github.com/stretchr/testify/assert"
)

func TestPi(t *testing.T) {
	val := gal.Pi()
	assert.Equal(t, gal.NewNumberFromFloat(math.Pi).String(), val.String())
}

func TestFactorial(t *testing.T) {
	val := gal.Factorial(gal.NewNumber(0))
	assert.Equal(t, gal.NewNumber(1).String(), val.String())

	val = gal.Factorial(gal.NewNumber(1))
	assert.Equal(t, gal.NewNumber(1).String(), val.String())

	val = gal.Factorial(gal.NewNumber(10))
	assert.Equal(t, gal.NewNumber(3_628_800).String(), val.String())
}
