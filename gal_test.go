package gal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	xpn := `-3 + 4`
	Eval(xpn)
}

func Test_extractPart1(t *testing.T) {
	expr := `-3 + -4`

	p, _, j := extractPart(expr)
	if !cmp.Equal("-3", p) {
		t.Error(cmp.Diff("-3", p))
	}
	i := j
	assert.Equal(t, numericalType, partType(p))

	p, _, j = extractPart(expr[i:])
	if !cmp.Equal("+", p) {
		t.Error(cmp.Diff("+", p))
	}
	i += j
	assert.Equal(t, operatorType, partType(p))

	p, _, j = extractPart(expr[i:])
	if !cmp.Equal("-4", p) {
		t.Error(cmp.Diff("-4", p))
	}
	i += j
	assert.Equal(t, numericalType, partType(p))

	assert.Equal(t, len(expr), i)
}

func Test_extractPart2(t *testing.T) {
	expr := `"-3 + -4" + -3 --4`

	p, _, j := extractPart(expr)
	if !cmp.Equal(`"-3 + -4"`, p) {
		t.Error(cmp.Diff(`"-3 + -4"`, p))
	}
	i := j
	assert.Equal(t, stringType, partType(p))

	p, _, j = extractPart(expr[i:])
	if !cmp.Equal("+", p) {
		t.Error(cmp.Diff("+", p))
	}
	i += j
	assert.Equal(t, operatorType, partType(p))

	p, _, j = extractPart(expr[i:])
	if !cmp.Equal("-3", p) {
		t.Error(cmp.Diff("-3", p))
	}
	i += j
	assert.Equal(t, numericalType, partType(p))

	p, _, j = extractPart(expr[i:])
	if !cmp.Equal("-", p) {
		t.Error(cmp.Diff("-", p))
	}
	i += j
	assert.Equal(t, operatorType, partType(p))

	p, _, j = extractPart(expr[i:])
	if !cmp.Equal("-4", p) {
		t.Error(cmp.Diff("-4", p))
	}
	i += j
	assert.Equal(t, numericalType, partType(p))

	assert.Equal(t, len(expr), i)
}

func TestParse(t *testing.T) {
	expr := `"-3 + -4" + -3 --4 / ( 1 + 2+3+4)`
	parseParts(expr)
}

func TestParse_Variable(t *testing.T) {
	expr := `:var_not_ended`
	parseParts(expr)

	expr = ":var with \nblanks:"
	parseParts(expr)

	expr = `:var_ended:`
	parseParts(expr)
}

func TestParse_FunctionName(t *testing.T) {
	expr := `f(abc)`
	parseParts(expr)

	expr = "f un c ti on   (" // Hmm, m'kay?
	parseParts(expr)

	expr = `func(`
	parseParts(expr)
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
