package gal_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/seborama/gal/v10"
)

func TestEval(t *testing.T) {
	expr := `-1 + 2 * 3 / 2 + 3 ** 2 -8`
	val := gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumberFromInt(3).String(), val.String())

	expr = `-"123"+"100"`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumberFromInt(-23).String(), val.String())

	expr = `1-2+7<<4+5`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumberFromInt((1-2+7)<<(4+5)).String(), val.String())

	expr = `-1-2-7<<4+5`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumberFromInt((-1-2-7)<<(4+5)).String(), val.String())

	expr = `-100*2*7+1>>2+3`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumberFromInt((-100*2*7+1)>>(2+3)).String(), val.String())

	expr = `100*2*7+1>>2+3`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumberFromInt((100*2*7+1)>>(2+3)).String(), val.String())

	expr = `2+Factorial(4)-5`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.NewNumberFromInt(21).String(), val.String())
}

func TestTreeBuilder_FromExpr_Variables(t *testing.T) {
	vars := gal.Variables{
		":var1:": gal.NewNumberFromInt(4),
		":var2:": gal.NewNumberFromInt(3),
	}

	expr := `2 + :var1: * :var2: - 5`

	got := gal.Parse(expr).Eval(gal.WithVariables(vars))
	expected := gal.NewNumberFromInt(9)

	if !cmp.Equal(expected, got) {
		t.Error(cmp.Diff(expected, got))
		t.FailNow()
	}
}

func TestTreeBuilder_FromExpr_UnknownVariable(t *testing.T) {
	expr := `2 + :var1: * :var2: - 5`

	got := gal.Parse(expr).Eval()
	expected := gal.NewUndefinedWithReasonf("error: unknown user-defined variable ':var1:'")

	if !cmp.Equal(expected, got) {
		t.Error(cmp.Diff(expected, got))
		t.FailNow()
	}
}

func TestEval_Boolean(t *testing.T) {
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

	// TODO: it should be easy to add support for "If" blocks.
	// ...   It could be in the form of am extension to the Bool type.
	// ...   Or, we could introduce a new built-in function "If" that returns a type that has methods Then() and Else()
	// ...   Or, it could be a built-in Object "If" that has methods Then() and Else()
	expr = `( 123 == 123 && 12 <= 45 ) Or ( "a" != "b" )`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	expr = `( 123 == 123 && 12 <= 45 ) Or ( "b" != "b" )`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	expr = `( 123 == 123 && 12 > 45 ) Or ( "b" == "b" )`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	expr = `( 123 == 123 And 12 > 45 ) Or ( "b" != "b" )`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.False.String(), val.String())

	expr = `True Or False`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	expr = `True Or (False)`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, gal.True.String(), val.String())

	// in this expression, the `()` are attached to `Or` which makes `Or()` a user-defined
	// function, rather than the `Or` operator.
	expr = `True Or(False)`
	val = gal.Parse(expr).Eval()
	assert.Equal(t, `undefined: error: unknown user-defined function 'Or'`, val.String())
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

			return value.Number().Multiply(gal.NewNumberFromInt(2))
		},
		"triple": func(args ...gal.Value) gal.Value {
			if len(args) != 1 {
				return gal.NewUndefinedWithReasonf("triple() requires a single argument, got %d", len(args))
			}

			value, ok := args[0].(gal.Numberer)
			if !ok {
				return gal.NewUndefinedWithReasonf("triple(): syntax error - argument must be a number-like value, got '%v'", args[0])
			}

			return value.Number().Multiply(gal.NewNumberFromInt(3))
		},
	}

	vars := gal.Variables{
		":val1:": gal.NewNumberFromInt(4),
		":val2:": gal.NewNumberFromInt(5),
	}

	got := parsedExpr.Eval(
		gal.WithVariables(vars),
		gal.WithFunctions(funcs),
	)
	expected := gal.NewNumberFromInt(23)

	if !cmp.Equal(expected, got) {
		t.Error(cmp.Diff(expected, got))
		t.FailNow()
	}

	// step 2: re-define funcs and vars and Eval the expression again
	// note that we do not need to parse the expression again, only just evaluate it
	funcs = gal.Functions{
		"double": func(args ...gal.Value) gal.Value {
			// should first validate argument count here
			value := args[0].(gal.Numberer) // should check type assertion is ok here
			return value.Number().Divide(gal.NewNumberFromInt(2))
		},
		"triple": func(args ...gal.Value) gal.Value {
			// should first validate argument count here
			value := args[0].(gal.Numberer) // should check type assertion is ok here
			return value.Number().Divide(gal.NewNumberFromInt(3))
		},
	}

	vars = gal.Variables{
		":val1:": gal.NewNumberFromInt(2),
		":val2:": gal.NewNumberFromInt(6),
	}

	got = parsedExpr.Eval(
		gal.WithVariables(vars),
		gal.WithFunctions(funcs),
	)
	expected = gal.NewNumberFromInt(3)

	if !cmp.Equal(expected, got) {
		t.Error(cmp.Diff(expected, got))
		t.FailNow()
	}
}

func TestNestedFunctions(t *testing.T) {
	expr := `double(triple(7))`
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

			return value.Number().Multiply(gal.NewNumberFromInt(2))
		},
		"triple": func(args ...gal.Value) gal.Value {
			if len(args) != 1 {
				return gal.NewUndefinedWithReasonf("triple() requires a single argument, got %d", len(args))
			}

			value, ok := args[0].(gal.Numberer)
			if !ok {
				return gal.NewUndefinedWithReasonf("triple(): syntax error - argument must be a number-like value, got '%v'", args[0])
			}

			return value.Number().Multiply(gal.NewNumberFromInt(3))
		},
	}

	got := parsedExpr.Eval(
		gal.WithFunctions(funcs),
	)
	expected := gal.NewNumberFromInt(42)
	assert.Equal(t, expected.String(), got.String())
}

// If renaming this test, also update the README.md file, where it is mentioned.
func TestMultiValueFunctions(t *testing.T) {
	expr := `sum(div(triple(7) double(4)))`

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

			return value.Number().Multiply(gal.NewNumberFromInt(2))
		},
		"triple": func(args ...gal.Value) gal.Value {
			if len(args) != 1 {
				return gal.NewUndefinedWithReasonf("triple() requires a single argument, got %d", len(args))
			}

			value, ok := args[0].(gal.Numberer)
			if !ok {
				return gal.NewUndefinedWithReasonf("triple(): syntax error - argument must be a number-like value, got '%v'", args[0])
			}

			return value.Number().Multiply(gal.NewNumberFromInt(3))
		},
		"div": func(args ...gal.Value) gal.Value {
			// returns the division of value1 by value2 as the interger portion and the remainder
			if len(args) != 2 {
				return gal.NewUndefinedWithReasonf("div() requires two arguments, got %d", len(args))
			}

			dividend := args[0].(gal.Numberer).Number()
			divisor := args[1].(gal.Numberer).Number()

			quotient := dividend.Divide(divisor).(gal.Numberer).Number().IntPart()
			remainder := dividend.Number().Sub(quotient.(gal.Number).Multiply(divisor.Number()))
			return gal.NewMultiValue(quotient, remainder)
		},
		"sum": func(args ...gal.Value) gal.Value {
			// NOTE: we convert the args to a MultiValue to make this function "bilingual".
			// That way, it can receiv either two Numberer's or one single MultiValue that holds 2 Numberer's.
			var margs gal.MultiValue
			if len(args) == 1 {
				margs = args[0].(gal.MultiValue) // not checking type satisfaction for simplicity
			}
			if len(args) == 2 {
				margs = gal.NewMultiValue(args...)
			}
			if margs.Size() != 2 {
				return gal.NewUndefinedWithReasonf("sum() requires either two Numberer-type Value's or one MultiValue holdings 2 Numberer's, as arguments, but got %d arguments", margs.Size())
			}

			value1 := args[0].(gal.MultiValue).Get(0).(gal.Numberer)
			value2 := args[0].(gal.MultiValue).Get(1).(gal.Numberer)

			return value1.Number().Add(value2.Number())
		},
	}

	got := gal.Parse(expr).Eval(
		gal.WithFunctions(funcs),
	)
	expected := gal.NewNumberFromInt(7)
	assert.Equal(t, expected.String(), got.String())
}

func TestStringsWithSpaces(t *testing.T) {
	expr := `"ab cd" + "ef gh"`
	parsedExpr := gal.Parse(expr)

	got := parsedExpr.Eval()
	assert.Equal(t, `"ab cdef gh"`, got.String())
}

func TestFunctionsAndStringsWithSpaces(t *testing.T) {
	expr := `f("ab cd") + f("ef gh")`
	parsedExpr := gal.Parse(expr)

	got := parsedExpr.Eval(
		gal.WithFunctions(gal.Functions{
			"f": func(args ...gal.Value) gal.Value {
				if len(args) != 1 {
					return gal.NewUndefinedWithReasonf("f() requires a single argument, got %d", len(args))
				}
				return args[0]
			},
		}),
	)
	assert.Equal(t, `"ab cdef gh"`, got.String())
}

func TestObjects_Properties(t *testing.T) {
	expr := `aCar.MaxSpeed - aCar.Speed`
	parsedExpr := gal.Parse(expr)

	expectedTree := gal.Tree{
		gal.NewObjectProperty("aCar", "MaxSpeed"),
		gal.Minus,
		gal.NewObjectProperty("aCar", "Speed"),
	}

	assert.Equal(t, expectedTree, parsedExpr)

	got := parsedExpr.Eval(
		gal.WithObjects(map[string]gal.Object{
			"aCar": Car{
				Make:     "Lotus Esprit",
				Mileage:  gal.NewNumberFromInt(2000),
				Speed:    100,
				MaxSpeed: 250,
			},
		}),
	)
	assert.Equal(t, "150", got.String())
}

func TestObjects_ChainedProperties(t *testing.T) {
	expr := `aCar.Stereo.Brand.Name + "::" + aCar.Stereo.Brand.Country`
	parsedExpr := gal.Parse(expr)

	expectedTree := gal.Tree{
		gal.NewObjectProperty("aCar", "Stereo"),
		gal.Dot[gal.Variable]{
			Member: gal.NewVariable(
				"Brand",
			),
		},
		gal.Dot[gal.Variable]{
			Member: gal.NewVariable(
				"Name",
			),
		},
		gal.Plus,
		gal.NewString("::"),
		gal.Plus,
		gal.NewObjectProperty("aCar", "Stereo"),
		gal.Dot[gal.Variable]{
			Member: gal.NewVariable(
				"Brand",
			),
		},
		gal.Dot[gal.Variable]{
			Member: gal.NewVariable(
				"Country",
			),
		},
	}

	require.Equal(t, expectedTree, parsedExpr)

	got := parsedExpr.Eval(
		gal.WithObjects(map[string]gal.Object{
			"aCar": Car{
				Make:     "Lotus Esprit",
				Mileage:  gal.NewNumberFromInt(2000),
				Speed:    100,
				MaxSpeed: 250,
				Stereo: CarStereo{
					Brand: StereoBrand{
						Name:    "Mitsubishi",
						Country: "Japan",
					},
					MaxWattage: 120,
				},
			},
		}),
	)
	assert.Equal(t, "Mitsubishi::Japan", got.AsString().RawString())
}

func TestObjects_Properties_TwoObjects(t *testing.T) {
	expr := `Road.Type == "Highway"
	And Car.IsRunning()
	And Car.Speed < 100
	And Car.Speed <= Car.MaxSpeed`
	parsedExpr := gal.Parse(expr)

	got := parsedExpr.Eval(
		gal.WithObjects(map[string]gal.Object{
			"Car": &Car{
				Make:     "Lotus Esprit",
				Mileage:  gal.NewNumberFromInt(2000),
				Speed:    80,
				MaxSpeed: 250,
			},
			"Road": &Road{
				Type: "Highway",
			},
		}),
	)
	assert.Equal(t, gal.True, got)
}

func TestObjects_Methods(t *testing.T) {
	expr := `aCar.MaxSpeed - aCar.CurrentSpeed()`
	parsedExpr := gal.Parse(expr)

	got := parsedExpr.Eval(
		gal.WithObjects(map[string]gal.Object{
			"aCar": &Car{
				Make:     "Lotus Esprit",
				Mileage:  gal.NewNumberFromInt(2000),
				Speed:    100,
				MaxSpeed: 250,
			},
		}),
	)
	assert.Equal(t, "150", got.String())
}

func TestObjects_Methods_WithSubTree(t *testing.T) {
	expr := `2 * (aCar.MaxSpeed - aCar.CurrentSpeed())`
	parsedExpr := gal.Parse(expr)

	got := parsedExpr.Eval(
		gal.WithObjects(map[string]gal.Object{
			"aCar": &Car{
				Make:     "Lotus Esprit",
				Mileage:  gal.NewNumberFromInt(2000),
				Speed:    100,
				MaxSpeed: 250,
			},
		}),
	)
	assert.Equal(t, "300", got.String())
}

func TestObjects_Methods_WithArgsSubTree(t *testing.T) {
	expr := `2 * (aCar.MaxSpeed - aCar.TillMaxSpeed(aCar.CurrentSpeed()))`
	parsedExpr := gal.Parse(expr)

	got := parsedExpr.Eval(
		gal.WithObjects(map[string]gal.Object{
			"aCar": &Car{
				Make:     "Lotus Esprit",
				Mileage:  gal.NewNumberFromInt(2000),
				Speed:    100,
				MaxSpeed: 250,
			},
		}),
	)
	assert.Equal(t, "200", got.String())
}

func TestObjects_MethodReceiver(t *testing.T) {
	expr := `aCar.MaxSpeed - aCar.CurrentSpeed()`
	parsedExpr := gal.Parse(expr)

	got := parsedExpr.Eval(
		gal.WithObjects(map[string]gal.Object{
			"aCar": Car{
				Make:     "Lotus Esprit",
				Mileage:  gal.NewNumberFromInt(2000),
				Speed:    100,
				MaxSpeed: 250,
			},
		}),
	)
	// Note: in this test, WithObjects is called with a `Car`, not a `*Car`.
	// However, Car.CurrentSpeed has a *Car receiver, hence from a Go perspective, the method
	// exists on *Car but it does NOT exist on Car!
	assert.Equal(t, "undefined: type 'gal_test.Car' does not have a method 'CurrentSpeed' (check if it has a pointer receiver)", got.String())
}

// TODO: this is an idea for a future feature.
// func TestObjects_Collections(t *testing.T) {
//	expr := `aCar.Tyres[0].Age - aCar.Drivers["Bob"].Age`
//	parsedExpr := gal.Parse(expr)

//	got := parsedExpr.Eval(
//		gal.WithObjects(map[string]gal.Object{
//			"aCar": &Car{
//				Make:            "Lotus Esprit",
//				Mileage:         gal.NewNumberFromInt(2000),
//				Speed:           100,
//				MaxSpeed:        250,
//				ComplexProperty: 0,
//				Tyres: []Tyre{
//					{
//						Location: "Front-Left",
//						Age:      5,
//					},
//					{
//						Location: "Front-Right",
//						Age:      1,
//					},
//				},
//				Drivers: map[string]Driver{
//					"Bob": {
//						Age: 32,
//					},
//					"Alice": {
//						Age: 37,
//					},
//				},
//			},
//		}),
//	)
//	assert.Equal(t, "-27", got.String())
// }
