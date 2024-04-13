package gal_test

import (
	"math"
	"testing"

	"github.com/seborama/gal/v7"
	"github.com/stretchr/testify/assert"
)

func TestPi(t *testing.T) {
	val := gal.Pi()
	assert.Equal(t, gal.NewNumberFromFloat(math.Pi).String(), val.String())
}

func TestPiLong(t *testing.T) {
	val := gal.PiLong()
	assert.Equal(t, math.Pi, gal.ToNumber(val).Float64())
}

func TestFactorial(t *testing.T) {
	val := gal.Factorial(gal.NewNumber(0))
	assert.Equal(t, gal.NewNumber(1).String(), val.String())

	val = gal.Factorial(gal.NewNumber(1))
	assert.Equal(t, gal.NewNumber(1).String(), val.String())

	val = gal.Factorial(gal.NewNumber(10))
	assert.Equal(t, gal.NewNumber(3_628_800).String(), val.String())
}

func TestCos(t *testing.T) {
	val := gal.Cos(gal.Pi())
	assert.Equal(t, -1.0, gal.ToNumber(val).Float64())
}

func TestSin(t *testing.T) {
	val := gal.Sin(gal.Pi().Divide(gal.NewNumberFromFloat(2.0)))
	assert.Equal(t, 1.0, gal.ToNumber(val).Float64())
}

func TestTan(t *testing.T) {
	val := gal.Tan(gal.NewNumberFromFloat(1.57079632))
	assert.Equal(t, gal.ToNumber(val).Int64(), int64(147169275))
}

func TestSqrt(t *testing.T) {
	val := gal.Sqrt(gal.NewNumberFromFloat(2.0))
	assert.InEpsilon(t, gal.ToNumber(val).Float64(), 1.414213562, 0.000001)
}

func TestFloor(t *testing.T) {
	val := gal.Floor(gal.NewNumberFromFloat(0.0))
	assert.Equal(t, int64(0), gal.ToNumber(val).Int64())
}
