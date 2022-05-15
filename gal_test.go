package gal_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/seborama/gal"
	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	xpn := `-3 + 4`
	val := gal.Parse(xpn).Eval()
	assert.Equal(t, gal.NewNumber(1), val)
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
	funcs := gal.Functions{
		"double": func(args ...gal.Value) gal.Value {
			if len(args) != 1 {
				return gal.NewUndefinedWithReasonf("double() requires a single argument, got %d", len(args))
			}

			v, ok := args[0].(gal.Numberer)
			if !ok {
				return gal.NewUndefinedWithReasonf("double(): syntax error - argument must be a number-like value, got '%v'", args[0])
			}

			return v.Number().Multiply(gal.NewNumber(2))
		},
		"triple": func(args ...gal.Value) gal.Value {
			if len(args) != 1 {
				return gal.NewUndefinedWithReasonf("triple() requires a single argument, got %d", len(args))
			}

			v, ok := args[0].(gal.Numberer)
			if !ok {
				return gal.NewUndefinedWithReasonf("triple(): syntax error - argument must be a number-like value, got '%v'", args[0])
			}

			return v.Number().Multiply(gal.NewNumber(3))
		},
	}

	vars := gal.Variables{
		":val1:": gal.NewNumber(4),
		":val2:": gal.NewNumber(5),
	}

	expr := `double(:val1:) + triple(:val2:)`

	got := gal.
		Parse(expr, gal.WithFunctions(funcs)).
		Eval(gal.WithVariables(vars))
	expected := gal.NewNumber(23)

	if !cmp.Equal(expected, got) {
		t.Error(cmp.Diff(expected, got))
	}
}
