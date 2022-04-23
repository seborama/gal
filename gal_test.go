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

func Test_buildExprTree_PlusMinus_String(t *testing.T) {
	expr := `"-3 + -4" + -3 --4 / ( 1 + 2+3+4) +log(10)`
	tree, err := buildExprTree(expr)
	require.NoError(t, err)

	expectedTree := Tree{
		NewString(`"-3 + -4"`),
		plus,
		minus,
		NewNumber(3),
		plus,
		NewNumber(4),
		dividedBy,
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

func Test_buildExprTree_VariousOperators(t *testing.T) {
	expr := `-1 + 2 * 3 / 2 + 1` // == 3  // -1 + ( 2 * 3 / 2 ) + 1
	tree, err := buildExprTree(expr)
	require.NoError(t, err)

	expectedTree := Tree{
		minus, // TODO: this should either be preceded with '0' or removed in favour of next line changing to NewNumber(-1)
		NewNumber(1),
		plus,
		NewNumber(2),
		times,
		NewNumber(3),
		dividedBy,
		NewNumber(2),
		plus,
		NewNumber(1),
	}

	if !cmp.Equal(expectedTree, tree) {
		t.Error(cmp.Diff(expectedTree, tree))
	}
}

func Test_prioritiseExprTreeOperators(t *testing.T) {
	inTree := Tree{
		NewString(`"-3 + -4"`),
		plus,
		minus,
		NewNumber(3),
		times,
		NewNumber(4),
		dividedBy,
		Tree{
			NewNumber(1),
			plus,
			NewNumber(2),
			plus,
			NewNumber(3),
			plus,
			NewNumber(4),
		},
		modulus,
		NewNumber(6),
		plus,
		NewFunction(
			"log",
			Tree{
				NewNumber(10),
			},
		),
	}
	// from: "-3 + -4" + -3 * 4 / (1 + 2 + 3 + 4) % 6 + log( 10 )
	//   to: "-3 + -4" + ( -3 * 4 / (1 + 2 + 3 + 4) % 6 ) + log( 10 )

	outTree := prioritiseExprTreeOperators(inTree)

	expectedTree := Tree{
		NewString(`"-3 + -4"`),
		plus,
		minus,
		Tree{
			NewNumber(3),
			times,
			NewNumber(4),
			dividedBy,
			Tree{
				NewNumber(1),
				plus,
				NewNumber(2),
				plus,
				NewNumber(3),
				plus,
				NewNumber(4),
			},
			modulus,
			NewNumber(6),
		},
		plus,
		NewFunction(
			"log",
			Tree{
				NewNumber(10),
			},
		),
	}

	if !cmp.Equal(expectedTree, outTree) {
		t.Error(cmp.Diff(expectedTree, outTree))
	}
}

func TestParse_Variable(t *testing.T) {
	expr := `:var_not_ended`
	_, err := parseParts(expr)
	require.Error(t, err)

	expr = ":var with \nblanks:"
	_, err = parseParts(expr)
	require.Error(t, err)

	expr = `:var_ended:`
	_, err = parseParts(expr)
	require.NoError(t, err)
}

func TestParse_FunctionName(t *testing.T) {
	expr := `f(4+g(5 6 (3+4))+ 6) + k() + (l(9))`
	_, err := parseParts(expr)
	require.NoError(t, err)

	expr = "f un c ti on   ("
	_, err = parseParts(expr)
	require.Error(t, err)

	expr = `func(`
	_, err = parseParts(expr)
	require.Error(t, err)
}

func TestParse_Operator(t *testing.T) {
	expr := `+++----+---+2`
	parseParts(expr)
}

func TestParse_Number(t *testing.T) {
	expr := `2`
	parseParts(expr)

	expr = `2.123`
	parseParts(expr)

	expr = `2.12.3`
	parseParts(expr)

	expr = `2.1 2.3`
	parseParts(expr)
}
