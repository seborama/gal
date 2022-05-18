package gal_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/seborama/gal"
	"github.com/stretchr/testify/require"
)

func TestTreeBuilder_FromExpr_VariousOperators(t *testing.T) {
	expr := `-1 + 2 * 3 / 2 + 1` // == 3  // -1 + ( 2 * 3 / 2 ) + 1
	tree, err := gal.NewTreeBuilder().FromExpr(expr)
	require.NoError(t, err)

	expectedTree := gal.Tree{
		gal.Minus,
		gal.NewNumber(1),
		gal.Plus,
		gal.NewNumber(2),
		gal.Multiply,
		gal.NewNumber(3),
		gal.Divide,
		gal.NewNumber(2),
		gal.Plus,
		gal.NewNumber(1),
	}

	if !cmp.Equal(expectedTree, tree) {
		t.Error(cmp.Diff(expectedTree, tree))
	}
}

func TestTreeBuilder_FromExpr_PlusMinus_String(t *testing.T) {
	expr := `"-3 + -4" + -3 --4 / ( 1 + 2+3+4) +tan(10)`
	tree, err := gal.NewTreeBuilder().FromExpr(expr)
	require.NoError(t, err)

	expectedTree := gal.Tree{
		gal.NewString(`-3 + -4`),
		gal.Minus,
		gal.NewNumber(3),
		gal.Plus,
		gal.NewNumber(4),
		gal.Divide,
		gal.Tree{
			gal.NewNumber(1),
			gal.Plus,
			gal.NewNumber(2),
			gal.Plus,
			gal.NewNumber(3),
			gal.Plus,
			gal.NewNumber(4),
		},
		gal.Plus,
		gal.NewFunction(
			"tan",
			gal.Tan,
			gal.Tree{
				gal.NewNumber(10),
			},
		),
	}

	if !cmp.Equal(expectedTree, tree) {
		t.Error(cmp.Diff(expectedTree, tree))
	}
}

func TestTreeBuilder_FromExpr_Functions(t *testing.T) {
	expr := `trunc(tan(10 + sin(cos(3 + f(1+2 3 ")4((")))) 6)`

	funcs := gal.Functions{
		"f": func(...gal.Value) gal.Value { return gal.NewNumber(123) },
	}

	got := gal.Parse(expr)

	expectedTree := gal.Tree{
		gal.NewFunction(
			"trunc",
			gal.Trunc,
			gal.Tree{
				gal.NewFunction(
					"tan",
					gal.Tan,
					gal.Tree{
						gal.NewNumber(10),
						gal.Plus,
						gal.NewFunction(
							"sin",
							gal.Sin,
							gal.Tree{
								gal.NewFunction(
									"cos",
									gal.Cos,
									gal.Tree{
										gal.NewNumber(3),
										gal.Plus,
										gal.NewFunction(
											"f",
											nil,
											gal.Tree{
												gal.NewNumber(1),
												gal.Plus,
												gal.NewNumber(2),
											},
											gal.Tree{
												gal.NewNumber(3),
											},
											gal.Tree{
												gal.NewString(")4(("),
											},
										),
									},
								),
							},
						),
					},
				),
			},
			gal.Tree{
				gal.NewNumber(6),
			},
		),
	}

	if !cmp.Equal(expectedTree, got) {
		t.Error(cmp.Diff(expectedTree, got))
		t.FailNow()
	}

	gotVal := got.Eval(gal.WithFunctions(funcs))
	expectedVal := gal.NewNumberFromFloat(5.323784)

	if !cmp.Equal(expectedVal, gotVal) {
		t.Error(cmp.Diff(expectedVal, gotVal))
	}
}
