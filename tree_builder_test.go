package gal_test

import (
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
		t.FailNow()
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
		t.FailNow()
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
		t.FailNow()
	}
}

func TestTreeBuilder_FromExpr_Objects(t *testing.T) {
	expr := `aCar.MaxSpeed - aCar.CurrentSpeed()`

	got := gal.Parse(expr)

	expectedTree := gal.Tree{
		gal.NewObjectProperty("aCar", "MaxSpeed"),
		gal.Minus,
		gal.NewObjectMethod("aCar", "CurrentSpeed", []gal.Tree{}...),
	}

	if !cmp.Equal(expectedTree, got) {
		t.Error(cmp.Diff(expectedTree, got))
		t.FailNow()
	}
}

func TestTreeBuilder_FromExpr_Dot_Accessor_Function(t *testing.T) {
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	// defer func() { slog.SetLogLoggerLevel(slog.LevelInfo) }()

	expr := `aCar.CurrentSpeed().Add(50).Add( 10 + (aCar.GetMaxSpeed() + aCar.MaxSpeed) ).Sub(20) - 100 - Factorial(5).Multiply(2)`

	got := gal.Parse(expr)

	expectedTree := gal.Tree{
		gal.NewObjectMethod("aCar", "CurrentSpeed", []gal.Tree{}...), // returns a "Number" which is a "Value"
		gal.Dot[gal.Function]{
			Member: gal.NewFunction(
				"Add",
				nil,
				gal.Tree{
					gal.NewNumberFromInt(50),
				},
			),
		},
		gal.Dot[gal.Function]{
			Member: gal.NewFunction(
				"Add",
				nil,
				gal.Tree{
					gal.NewNumberFromInt(10),
					gal.Plus,
					gal.Tree{
						gal.NewObjectMethod("aCar", "GetMaxSpeed", []gal.Tree{}...),
						gal.Plus,
						gal.NewObjectProperty("aCar", "MaxSpeed"),
					},
				},
			),
		},
		gal.Dot[gal.Function]{
			Member: gal.NewFunction(
				"Sub",
				nil,
				gal.Tree{
					gal.NewNumberFromInt(20),
				},
			),
		},
		gal.Minus,
		gal.NewNumber(100, 0),
		gal.Minus,
		gal.Function{
			Name:   "Factorial",
			BodyFn: gal.Factorial,
			Args: []gal.Tree{
				{gal.NewNumberFromInt(5)},
			},
		},
		gal.Dot[gal.Function]{
			Member: gal.NewFunction(
				"Multiply",
				nil,
				gal.Tree{
					gal.NewNumberFromInt(2),
				},
			),
		},
	}

	if !cmp.Equal(expectedTree, got) {
		t.Error(cmp.Diff(expectedTree, got))
		t.FailNow()
	}

	gotVal := got.Eval(
		gal.WithObjects(map[string]gal.Object{
			"aCar": &Car{
				Speed:    80,
				MaxSpeed: 200,
			},
		}),
	)
	assert.Equal(t, gal.NewNumberFromInt(180), gotVal)
}

func TestTreeBuilder_FromExpr_Dot_Accessor_Property(t *testing.T) {
	expr := `aCar.CurrentSpeed3().Speed`

	got := gal.Parse(expr)

	expectedTree := gal.Tree{
		gal.NewObjectMethod("aCar", "CurrentSpeed3", []gal.Tree{}...), // returns a "fancyType"
		gal.Dot[gal.Variable]{
			Member: gal.NewVariable(
				"Speed",
			),
		},
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
	assert.Equal(t, gal.NewNumberFromFloat(100), gotVal)
}

// TODO: (!) this is an idea for a future feature
// func TestTreeBuilder_FromExpr_Arrays(t *testing.T) {
//	expr := `f(1 2 3)[1]` // simple example
//	expr := `f(1 2 3)[1 + sum(1 2 3)]` // complex example

//	funcs := gal.Functions{
//		"f": func(args ...gal.Value) gal.Value { return gal.NewMultiValue(args...) },
//	}

//	got := gal.Parse(expr)

//	expectedTree := gal.Tree{
//		gal.NewFunction(
//			"f",
//			nil,
//			gal.Tree{
//				gal.NewNumberFromInt(1),
//			},
//			gal.Tree{
//				gal.NewNumberFromInt(2),
//			},
//			gal.Tree{
//				gal.NewNumberFromInt(3),
//			},
//			gal.Collection{ // this would work similarly to gal.Dot but with the lhs as the collection and the rhs as the index
//				gal.Tree{
//					gal.NewNumberFromInt(1),
//				},
//			},
//		),
//	}

//	if !cmp.Equal(expectedTree, got) {
//		t.Error(cmp.Diff(expectedTree, got))
//      t.FailNow()
//	}

//	gotVal := got.Eval(gal.WithFunctions(funcs))
//	expectedVal := gal.NewNumberFromFloat(5.323784)

//	if !cmp.Equal(expectedVal, gotVal) {
//		t.Error(cmp.Diff(expectedVal, gotVal))
//      t.FailNow()
//	}
// }
