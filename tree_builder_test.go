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
	expr := `"-3 + -4" + -3 --4 / ( 1 + 2+3+4) +log(10)`
	tree, err := gal.NewTreeBuilder().FromExpr(expr)
	require.NoError(t, err)

	expectedTree := gal.Tree{
		gal.NewString(`"-3 + -4"`),
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
			"log",
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
	expr := `log(10 + sin(cos(3 + f(1+2 3 4))))`
	tree, err := gal.NewTreeBuilder().FromExpr(expr)
	require.NoError(t, err)

	expectedTree := gal.Tree{
		gal.NewFunction(
			"log",
			gal.Tree{
				gal.NewNumber(10),
				gal.Plus,
				gal.NewFunction(
					"sin",
					gal.Tree{
						gal.NewFunction(
							"cos",
							gal.Tree{
								gal.NewNumber(3),
								gal.Plus,
								gal.NewFunction(
									"f",
									gal.Tree{
										gal.NewNumber(1),
										gal.Plus,
										gal.NewNumber(2),
									},
									gal.Tree{
										gal.NewNumber(3),
									},
									gal.Tree{
										gal.NewNumber(4),
									},
								),
							},
						),
					},
				),
			},
		),
	}

	if !cmp.Equal(expectedTree, tree) {
		t.Error(cmp.Diff(expectedTree, tree))
	}
}

func TestTreeBuilder_FromExpr_Variables(t *testing.T) {
	vars := map[string]gal.Value{
		":var1:": gal.NewNumber(4),
		":var2:": gal.NewNumber(3),
	}

	expr := `2 + :var1: * :var2: - 5`

	tree, err := gal.NewTreeBuilder(gal.WithVariables(vars)).FromExpr(expr)
	require.NoError(t, err)

	got := tree.Eval()
	expected := gal.NewNumber(9)

	if !cmp.Equal(expected, got) {
		t.Error(cmp.Diff(expected, got))
	}
}
