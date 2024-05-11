package gal_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/seborama/gal/v9"
)

func TestTreeBuilder_FromExpr_VariousOperators(t *testing.T) {
	expr := `-10 + 2 * 7 / 2 + 5 ** 4 -8`
	tree, err := gal.NewTreeBuilder().FromExpr(expr)
	require.NoError(t, err)

	expectedTree := gal.Tree{
		gal.NewNumberFromInt(-1),
		gal.Multiply,
		gal.NewNumberFromInt(10),
		gal.Plus,
		gal.NewNumberFromInt(2),
		gal.Multiply,
		gal.NewNumberFromInt(7),
		gal.Divide,
		gal.NewNumberFromInt(2),
		gal.Plus,
		gal.NewNumberFromInt(5),
		gal.Power,
		gal.NewNumberFromInt(4),
		gal.Minus,
		gal.NewNumberFromInt(8),
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
		gal.NewNumberFromInt(3),
		gal.Plus,
		gal.NewNumberFromInt(4),
		gal.Divide,
		gal.Tree{
			gal.NewNumberFromInt(1),
			gal.Plus,
			gal.NewNumberFromInt(2),
			gal.Plus,
			gal.NewNumberFromInt(3),
			gal.Plus,
			gal.NewNumberFromInt(4),
		},
		gal.Plus,
		gal.NewFunction(
			"tan",
			gal.Tan,
			gal.Tree{
				gal.NewNumberFromInt(10),
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
		"f": func(...gal.Value) gal.Value { return gal.NewNumberFromInt(123) },
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
						gal.NewNumberFromInt(10),
						gal.Plus,
						gal.NewFunction(
							"sin",
							gal.Sin,
							gal.Tree{
								gal.NewFunction(
									"cos",
									gal.Cos,
									gal.Tree{
										gal.NewNumberFromInt(3),
										gal.Plus,
										gal.NewFunction(
											"f",
											nil,
											gal.Tree{
												gal.NewNumberFromInt(1),
												gal.Plus,
												gal.NewNumberFromInt(2),
											},
											gal.Tree{
												gal.NewNumberFromInt(3),
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
				gal.NewNumberFromInt(6),
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
