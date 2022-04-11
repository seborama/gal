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

func Test_extractPart(t *testing.T) {
	expr := `-3 + -4`

	p, j := extractPart(expr)
	if !cmp.Equal("-3", p) {
		t.Error(cmp.Diff("-3", p))
	}
	i := j
	assert.Equal(t, numericalType, partType(p))

	p, j = extractPart(expr[i:])
	if !cmp.Equal("+", p) {
		t.Error(cmp.Diff("+", p))
	}
	i += j
	assert.Equal(t, operatorType, partType(p))

	p, j = extractPart(expr[i:])
	if !cmp.Equal("-4", p) {
		t.Error(cmp.Diff("-4", p))
	}
	i += j
	assert.Equal(t, numericalType, partType(p))

	assert.Equal(t, len(expr), i)
}
