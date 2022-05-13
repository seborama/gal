package gal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_extractPart_VariableName(t *testing.T) {
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

func Test_extractPart_FunctionName(t *testing.T) {
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

func Test_extractPart_Operator(t *testing.T) {
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

func Test_extractPart_Number(t *testing.T) {
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
