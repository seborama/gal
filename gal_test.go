package gal_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/seborama/gal/v7"
	"github.com/stretchr/testify/assert"
)

func TestBoolean(t *testing.T) {
	expr := `2 > 1`
	val := gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	expr = `2 > 2`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.False.String(), val.String())

	expr = `2 >= 2`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	expr = `2 < 1`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.False.String(), val.String())

	expr = `2 < 2`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.False.String(), val.String())

	expr = `2 <= 2`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	expr = `2 != 2`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.False.String(), val.String())

	expr = `1 != 2`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	expr = `3 != 2`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	expr = `2 == 2`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	expr = `1 == 2`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.False.String(), val.String())

	expr = `3 == 2`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.False.String(), val.String())
}

func TestEval(t *testing.T) {
	expr := `-1 + 2 * 3 / 2 + 3 ** 2 -8`
	val := gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumber(3).String(), val.String())

	expr = `-"123"+"100"`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumber(-23).String(), val.String())

	expr = `1-2+7<<4+5`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumber((1-2+7)<<(4+5)).String(), val.String())

	expr = `-1-2-7<<4+5`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumber((-1-2-7)<<(4+5)).String(), val.String())

	expr = `-100*2*7+1>>2+3`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumber((-100*2*7+1)>>(2+3)).String(), val.String())

	expr = `100*2*7+1>>2+3`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumber((100*2*7+1)>>(2+3)).String(), val.String())

	expr = `2+Factorial(4)-5`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumber(21).String(), val.String())
}

func TestTreeBuilder_FromExpr_Variables(t *testing.T) {
	vars := gal.Variables{
		":var1:": gal.NewNumber(4), // TODO: remove the need to surround with `:`?
		":var2:": gal.NewNumber(3),
	}

	expr := `2 + :var1: * :var2: - 5`

	got := gal.Parse(expr).Eval(gal.WithVariables(vars))
	expected := gal.NewNumber(9)

	if !cmp.Equal(expected, got) {
		t.Error(cmp.Diff(expected, got))
	}
}

func TestTreeBuilder_FromExpr_UnknownVariable(t *testing.T) {
	expr := `2 + :var1: * :var2: - 5`

	got := gal.Parse(expr).Eval()
	expected := gal.NewUndefinedWithReasonf("syntax error: unknown variable name: ':var1:'")

	if !cmp.Equal(expected, got) {
		t.Error(cmp.Diff(expected, got))
	}
}

func TestWithVariablesAndFunctions(t *testing.T) {
	expr := `double(:val1:) + triple(:val2:)`
	parsedExpr := gal.Parse(expr)

	// step 1: define funcs and vars and Eval the expression
	funcs := gal.Functions{
		"double": func(args ...gal.Value) gal.Value {
			if len(args) != 1 {
				return gal.NewUndefinedWithReasonf("double() requires a single argument, got %d", len(args))
			}

			value, ok := args[0].(gal.Numberer)
			if !ok {
				return gal.NewUndefinedWithReasonf("double(): syntax error - argument must be a number-like value, got '%v'", args[0])
			}

			return value.Number().Multiply(gal.NewNumber(2))
		},
		"triple": func(args ...gal.Value) gal.Value {
			if len(args) != 1 {
				return gal.NewUndefinedWithReasonf("triple() requires a single argument, got %d", len(args))
			}

			value, ok := args[0].(gal.Numberer)
			if !ok {
				return gal.NewUndefinedWithReasonf("triple(): syntax error - argument must be a number-like value, got '%v'", args[0])
			}

			return value.Number().Multiply(gal.NewNumber(3))
		},
	}

	vars := gal.Variables{
		":val1:": gal.NewNumber(4),
		":val2:": gal.NewNumber(5),
	}

	got := parsedExpr.Eval(
		gal.WithVariables(vars),
		gal.WithFunctions(funcs),
	)
	expected := gal.NewNumber(23)

	if !cmp.Equal(expected, got) {
		t.Error(cmp.Diff(expected, got))
	}

	// step 2: re-define funcs and vars and Eval the expression again
	// note that we do not need to parse the expression again, only just evaluate it
	funcs = gal.Functions{
		"double": func(args ...gal.Value) gal.Value {
			// should first validate argument count here
			value := args[0].(gal.Numberer) // should check type assertion is ok here
			return value.Number().Divide(gal.NewNumber(2))
		},
		"triple": func(args ...gal.Value) gal.Value {
			// should first validate argument count here
			value := args[0].(gal.Numberer) // should check type assertion is ok here
			return value.Number().Divide(gal.NewNumber(3))
		},
	}

	vars = gal.Variables{
		":val1:": gal.NewNumber(2),
		":val2:": gal.NewNumber(6),
	}

	got = parsedExpr.Eval(
		gal.WithVariables(vars),
		gal.WithFunctions(funcs),
	)
	expected = gal.NewNumber(3)

	if !cmp.Equal(expected, got) {
		t.Error(cmp.Diff(expected, got))
	}
}
