package gal_test

import (
	"math"
	"testing"

	"github.com/seborama/gal/v10"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	val := gal.Factorial(gal.NewNumberFromInt(0))
	assert.Equal(t, gal.NewNumberFromInt(1).String(), val.String())

	val = gal.Factorial(gal.NewNumberFromInt(1))
	assert.Equal(t, gal.NewNumberFromInt(1).String(), val.String())

	val = gal.Factorial(gal.NewNumberFromInt(10))
	assert.Equal(t, gal.NewNumberFromInt(3_628_800).String(), val.String())

	val = gal.Factorial(gal.NewNumberFromInt(-10))
	assert.Equal(t, "undefined: Factorial: requires a positive integer, cannot accept -10", val.String())
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

func TestLn(t *testing.T) {
	val := gal.Ln(gal.NewNumberFromInt(123456), gal.NewNumberFromInt(5))
	assert.Equal(t, "11.72364", val.String())

	val = gal.Ln(gal.NewNumberFromInt(-123456), gal.NewNumberFromInt(5))
	assert.Equal(t, "undefined: Ln:cannot calculate natural logarithm for negative decimals", val.String())
}

func TestLog(t *testing.T) {
	val := gal.Log(gal.NewNumberFromInt(123456), gal.NewNumberFromInt(5))
	assert.Equal(t, "5.09151", val.String())

	val = gal.Log(gal.NewNumberFromInt(-123456), gal.NewNumberFromInt(5))
	assert.Equal(t, "undefined: Log:cannot calculate natural logarithm for negative decimals", val.String())

	val = gal.Log(gal.NewNumberFromInt(10_000_000), gal.NewNumberFromInt(0))
	assert.Equal(t, "7", val.String())
}

func TestFunctionEval(t *testing.T) {
	expr := `eval("7+22")*2`
	tree, err := gal.NewTreeBuilder().FromExpr(expr)
	require.NoError(t, err)

	assert.Equal(t, "58", tree.Eval().String())
}
