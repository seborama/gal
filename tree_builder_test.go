package gal_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/seborama/gal/v10"
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

func TestTreeBuilder_FromExpr_Objects(t *testing.T) {
	expr := `aCar.MaxSpeed - aCar.CurrentSpeed()`

	got := gal.Parse(expr)

	expectedTree := gal.Tree{
		gal.NewVariable(
			"aCar.MaxSpeed",
		),
		gal.Minus,
		gal.Function{
			"aCar.CurrentSpeed",
			nil,
			[]gal.Tree{},
		},
	}

	if !cmp.Equal(expectedTree, got) {
		t.Error(cmp.Diff(expectedTree, got))
		t.FailNow()
	}
}

func TestTreeBuilder_FromExpr_Dot_Accessor_Function(t *testing.T) {
	expr := `aCar.CurrentSpeed().Add(50)+100`

	got := gal.Parse(expr)

	expectedTree := gal.Tree{
		// TODO: TBC
	}

	if !cmp.Equal(expectedTree, got) {
		t.Error(cmp.Diff(expectedTree, got))
		t.FailNow()
	}

	gotVal := got.Eval(
		gal.WithObjects(map[string]gal.Object{
			"aCar": &Car{
				Speed: 80,
			},
		}),
	)
	assert.Equal(t, gal.NewNumberFromInt(230), gotVal)
}

func TestTreeBuilder_FromExpr_Dot_Accessor_Property(t *testing.T) {
	expr := `aCar.CurrentSpeed3().Speed`

	got := gal.Parse(expr)

	fmt.Printf("%#v\n", got)
	expectedTree := gal.Tree{
		//TODO: TBC
	}

	if !cmp.Equal(expectedTree, got) {
		t.Error(cmp.Diff(expectedTree, got))
		t.FailNow()
	}

	gotVal := got.Eval(
		gal.WithObjects(map[string]gal.Object{
			"aCar": &Car{
				Speed: 100,
			},
		}),
	)
	assert.Equal(t, gal.NewNumberFromInt(100), gotVal)
}

func TestTreeBuilder_FromExpr_Arrays(t *testing.T) {
	expr := `f(1 2 3)[1]`

	funcs := gal.Functions{
		"f": func(args ...gal.Value) gal.Value { return gal.NewMultiValue(args...) },
	}

	got := gal.Parse(expr)

	expectedTree := gal.Tree{
		gal.NewFunction(
			"f",
			nil,
			gal.Tree{
				gal.NewNumberFromInt(1),
			},
			gal.Tree{
				gal.NewNumberFromInt(2),
			},
			gal.Tree{
				gal.NewNumberFromInt(3),
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
