package gal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiValueString(t *testing.T) {
	v := NewMultiValue(NewNumberFromInt(123), NewString("abc"), NewBool(true))
	assert.Equal(t, `123,"abc",True`, v.String())
}
