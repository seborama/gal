package gal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	xpn := `-3 + 4`
	val := Eval(xpn)
	assert.Equal(t, NewNumber(1), val)
}

func TestParse_Variable(t *testing.T) {
	expr := `:var_not_ended`
	_, _, _, err := extractPart(expr)
	require.Error(t, err)

	expr = ":var with \nblanks:"
	_, _, _, err = extractPart(expr)
	require.Error(t, err)

	expr = `:var_ended:`
	s, et, i, err := extractPart(expr)
	require.NoError(t, err)
	assert.Equal(t, ":var_ended:", s)
	assert.Equal(t, variableType, et)
	assert.Equal(t, 11, i)
}

func TestParse_FunctionName(t *testing.T) {
	expr := `f(4+g(5 6 (3+4))+ 6) + k() + (l(9))`
	s, et, i, err := extractPart(expr)
	require.NoError(t, err)
	assert.Equal(t, "f(4+g(5 6 (3+4))+ 6)", s)
	assert.Equal(t, functionType, et)
	assert.Equal(t, 20, i)

	expr = `(4+g(5 6 (3+4))+ 6) + k() + (l(9))`
	s, et, i, err = extractPart(expr)
	require.NoError(t, err)
	assert.Equal(t, "(4+g(5 6 (3+4))+ 6)", s)
	assert.Equal(t, functionType, et)
	assert.Equal(t, 19, i)

	expr = "f un c ti on   ("
	_, _, _, err = extractPart(expr)
	require.Error(t, err)

	expr = `func(`
	_, _, _, err = extractPart(expr)
	require.Error(t, err)
}

func TestParse_Operator(t *testing.T) {
	expr := `+ ++- ---+-- -+2`
	s, et, i, err := extractPart(expr)
	require.NoError(t, err)
	assert.Equal(t, "-", s)
	assert.Equal(t, operatorType, et)
	assert.Equal(t, 15, i)

	expr = `+ ++- ---+-- --+2`
	s, et, i, err = extractPart(expr)
	require.NoError(t, err)
	assert.Equal(t, "+", s)
	assert.Equal(t, operatorType, et)
	assert.Equal(t, 16, i)
}

func TestParse_Number(t *testing.T) {
	expr := `2`
	s, et, i, err := extractPart(expr)
	require.NoError(t, err)
	assert.Equal(t, "2", s)
	assert.Equal(t, numericalType, et)
	assert.Equal(t, 1, i)

	expr = `2.123`
	s, et, i, err = extractPart(expr)
	require.NoError(t, err)
	assert.Equal(t, "2.123", s)
	assert.Equal(t, numericalType, et)
	assert.Equal(t, 5, i)

	expr = `2.12.3`
	_, _, _, err = extractPart(expr)
	require.EqualError(t, err, "syntax error: invalid character '.' for number '2.12.'")

	expr = `2.1 2.3`
	s, et, i, err = extractPart(expr)
	require.NoError(t, err)
	assert.Equal(t, "2.1", s)
	assert.Equal(t, numericalType, et)
	assert.Equal(t, 3, i)
}

func Test_buildExprTree_VariousOperators(t *testing.T) {
	expr := `-1 + 2 * 3 / 2 + 1` // == 3  // -1 + ( 2 * 3 / 2 ) + 1
	tree, err := buildExprTree(expr)
	require.NoError(t, err)

	expectedTree := Tree{
		minus,
		NewNumber(1),
		plus,
		NewNumber(2),
		multiply,
		NewNumber(3),
		divide,
		NewNumber(2),
		plus,
		NewNumber(1),
	}

	if !cmp.Equal(expectedTree, tree) {
		t.Error(cmp.Diff(expectedTree, tree))
	}
}

func Test_buildExprTree_PlusMinus_String(t *testing.T) {
	expr := `"-3 + -4" + -3 --4 / ( 1 + 2+3+4) +log(10)`
	tree, err := buildExprTree(expr)
	require.NoError(t, err)

	expectedTree := Tree{
		NewString(`"-3 + -4"`),
		minus,
		NewNumber(3),
		plus,
		NewNumber(4),
		divide,
		Tree{
			NewNumber(1),
			plus,
			NewNumber(2),
			plus,
			NewNumber(3),
			plus,
			NewNumber(4),
		},
		plus,
		NewFunction(
			"log",
			Tree{
				NewNumber(10),
			},
		),
	}

	if !cmp.Equal(expectedTree, tree) {
		t.Error(cmp.Diff(expectedTree, tree))
	}
}

func Test_buildExprTree_Functions(t *testing.T) {
	expr := `log(10 + sin(cos(3 + f(1 2 3))))`
	tree, err := buildExprTree(expr)
	require.NoError(t, err)

	expectedTree := Tree{
		NewFunction(
			"log",
			Tree{
				NewNumber(10),
				plus,
				NewFunction(
					"sin",
					Tree{
						NewFunction(
							"cos",
							Tree{
								NewNumber(3),
								plus,
								NewFunction(
									"f",
									Tree{
										// TODO: this may change as multi-arg functions aren't really supported currently.
										NewNumber(1),
										NewNumber(2),
										NewNumber(3),
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
