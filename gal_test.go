package gal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	xpn := `-3 + 4`
	Eval(xpn)
}

func Test_extractPart(t *testing.T) {
	expr := `"-3 + -4" + -3 --4`

	p, _, j, err := extractPart(expr)
	require.NoError(t, err)
	if !cmp.Equal(`"-3 + -4"`, p) {
		t.Error(cmp.Diff(`"-3 + -4"`, p))
	}
	i := j
	assert.Equal(t, stringType, partType(p))

	p, _, j, err = extractPart(expr[i:])
	require.NoError(t, err)
	if !cmp.Equal("+", p) {
		t.Error(cmp.Diff("+", p))
	}
	i += j
	assert.Equal(t, operatorType, partType(p))

	p, _, j, err = extractPart(expr[i:])
	require.NoError(t, err)
	if !cmp.Equal("-3", p) {
		t.Error(cmp.Diff("-3", p))
	}
	i += j
	assert.Equal(t, numericalType, partType(p))

	p, _, j, err = extractPart(expr[i:])
	require.NoError(t, err)
	if !cmp.Equal("-", p) {
		t.Error(cmp.Diff("-", p))
	}
	i += j
	assert.Equal(t, operatorType, partType(p))

	p, _, j, err = extractPart(expr[i:])
	require.NoError(t, err)
	if !cmp.Equal("-4", p) {
		t.Error(cmp.Diff("-4", p))
	}
	i += j
	assert.Equal(t, numericalType, partType(p))

	assert.Equal(t, len(expr), i)
}

func TestParse(t *testing.T) {
	// expr := `"-3 + -4" + -3 --4 / ( 1 + 2+3+4)`
	// v, err := parseParts(expr)
	// require.NoError(t, err)
	// if !cmp.Equal(Undefined{}, v) {
	// 	t.Error(cmp.Diff(Undefined{}, v))
	// }

	// expr := `-3 --4 / ( 1 + 2+3+4)`
	// v, err := parseParts(expr)
	// require.NoError(t, err)
	// if !cmp.Equal(NewNumberFromFloat(-2.6), v) {
	// 	t.Error(cmp.Diff(Undefined{}, v))
	// }

	expr := `-1 + 2 * 3 / 2 + 1` // == 3  // -1 + ( 2 * 3 ) / 2 + ( 1 )
	v, err := parseParts(expr)
	require.NoError(t, err)
	if !cmp.Equal(NewNumberFromFloat(-2.6), v) {
		t.Error(cmp.Diff(Undefined{}, v))
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
