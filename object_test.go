package gal_test

import (
	"fmt"
	"testing"

	"github.com/seborama/gal/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Car struct {
	Make     string
	Mileage  gal.Number
	Speed    float32
	MaxSpeed int64
}

func (c *Car) String() string {
	return fmt.Sprint("Car", c.Make, c.Mileage.String(), c.Speed, c.Make)
}

func TestObjectGetProperty(t *testing.T) {
	myCar := &Car{
		Make:     "Lotus",
		Mileage:  gal.NewNumberFromInt(100),
		Speed:    50.345,
		MaxSpeed: 250,
	}

	var nilCar *Car

	val, ok := gal.ObjectGetProperty(nilCar, "Mileage")
	require.False(t, ok)
	assert.Equal(t, "undefined: object is nil, not a Go value or invalid", val.String())

	val, ok = gal.ObjectGetProperty(myCar, "Make")
	require.True(t, ok)
	assert.Equal(t, "Lotus", val.(gal.String).RawString())

	val, ok = gal.ObjectGetProperty(myCar, "Mileage")
	require.True(t, ok)
	assert.Equal(t, gal.NewNumberFromInt(100), val)

	// some bothersome floating point issues...
	val, ok = gal.ObjectGetProperty(*myCar, "Speed")
	require.True(t, ok)
	assert.Equal(t, gal.NewNumberFromFloat(50.345).Trunc(1).String(), val.(gal.Number).Trunc(1).String())

	val, ok = gal.ObjectGetProperty(complex(10, 3), "Blah")
	require.False(t, ok)
	assert.Equal(t, "undefined: object is 'complex128' but only 'struct' and '*struct' are currently supported", val.String())

	val, ok = gal.ObjectGetProperty(myCar, "MaxSpeed")
	require.True(t, ok)
	assert.Equal(t, gal.NewNumberFromInt(250), val)
}
