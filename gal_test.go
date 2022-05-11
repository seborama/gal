package gal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	xpn := `-3 + 4`
	Eval(xpn)
}

func TestParse_Variable(t *testing.T) {
	expr := `:var_not_ended`
	_, _, _, err := extractPart(expr)
	require.Error(t, err)

	expr = ":var with \nblanks:"
	_, _, _, err = extractPart(expr)
	require.Error(t, err)

	expr = `:var_ended:`
	_, _, _, err = extractPart(expr)
	require.NoError(t, err)
}

func TestParse_FunctionName(t *testing.T) {
	expr := `f(4+g(5 6 (3+4))+ 6) + k() + (l(9))`
	_, _, _, err := extractPart(expr)
	require.NoError(t, err)

	expr = "f un c ti on   ("
	_, _, _, err = extractPart(expr)
	require.Error(t, err)

	expr = `func(`
	_, _, _, err = extractPart(expr)
	require.Error(t, err)
}

func TestParse_Operator(t *testing.T) {
	expr := `+ ++- ---+-- -+2`
	extractPart(expr)
}

func TestParse_Number(t *testing.T) {
	expr := `2`
	_, _, _, err := extractPart(expr)
	require.NoError(t, err)

	expr = `2.123`
	_, _, _, err = extractPart(expr)
	require.NoError(t, err)

	expr = `2.12.3`
	_, _, _, err = extractPart(expr)
	require.EqualError(t, err, "syntax error: invalid character '.' for number '2.12.'")

	expr = `2.1 2.3`
	_, _, _, err = extractPart(expr)
	require.NoError(t, err)
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
