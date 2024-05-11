package gal_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/seborama/gal/v8"
)

type Car struct {
	Make            string
	Mileage         gal.Number
	Speed           float32
	MaxSpeed        int64
	ComplexProperty complex128
}

func (c *Car) Ignite() gal.Value {
	return gal.True
}

func (c Car) Shutdown() gal.Value { //nolint: gocritic
	return gal.True
}

func (c *Car) CurrentSpeed() gal.Value {
	return gal.NewNumberFromFloat(float64(c.Speed))
}

func (c *Car) CurrentSpeed2() float32 {
	return c.Speed
}

func (c *Car) SetSpeed(speed gal.Number) {
	c.Speed = float32(speed.Float64())
}

func (c *Car) SetSpeed2(speed gal.Number) gal.Bool {
	c.Speed = float32(speed.Float64())
	return gal.True
}

func (c *Car) SetSpeed3(speed float32) gal.Bool {
	c.Speed = speed
	return gal.True
}

func (c *Car) TillMaxSpeed(speed gal.Number) gal.Number {
	return gal.NewNumberFromInt(c.MaxSpeed).Add(speed.Neg()).(gal.Numberer).Number()
}

type fancyType struct {
	Speed float32
}

func (c *Car) SetSpeed4(speed fancyType) gal.Bool {
	c.Speed = speed.Speed
	return gal.True
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

	val, ok = gal.ObjectGetProperty(myCar, "ComplexProperty")
	require.False(t, ok)
	assert.Equal(t, "undefined: object::*gal_test.Car:ComplexProperty - type 'complex128' cannot be mapped to gal.Value", val.String())

	val, ok = gal.ObjectGetProperty(myCar, "MaxSpeed")
	require.True(t, ok)
	assert.Equal(t, gal.NewNumberFromInt(250), val)
}

func TestObjectGetMethod(t *testing.T) {
	myCar := &Car{
		Make:     "Lotus",
		Mileage:  gal.NewNumberFromInt(100),
		Speed:    50.345,
		MaxSpeed: 250,
	}

	var nilCar *Car

	val, ok := gal.ObjectGetMethod(nilCar, "Ignite")
	require.False(t, ok)
	assert.Equal(t, "undefined: object is nil for type '*gal_test.Car' or does not have a method 'Ignite'", val().String())

	val, ok = gal.ObjectGetMethod(myCar, "DoesNotExist!")
	require.False(t, ok)
	assert.Equal(t, "undefined: type '*gal_test.Car' does not have a method 'DoesNotExist!' (check if it has a pointer receiver)", val().String())

	val, ok = gal.ObjectGetMethod(myCar, "Ignite")
	require.True(t, ok)
	assert.Equal(t, gal.True, val())

	val, ok = gal.ObjectGetMethod(myCar, "CurrentSpeed")
	require.True(t, ok)
	assert.Equal(t, "50.345", val().(gal.Numberer).Number().Trunc(3).String())

	val, ok = gal.ObjectGetMethod(myCar, "CurrentSpeed2")
	require.True(t, ok)
	assert.Equal(t, "50.345", val().(gal.Numberer).Number().Trunc(3).String())

	val, ok = gal.ObjectGetMethod(myCar, "SetSpeed")
	require.True(t, ok)
	got := val(gal.NewNumberFromFloat(76), gal.True, gal.False)
	assert.Equal(t, "undefined: invalid function call - object::*gal_test.Car:SetSpeed - wants 1 args, received 3 instead", got.String())

	val, ok = gal.ObjectGetMethod(myCar, "SetSpeed")
	require.True(t, ok)
	got = val(gal.NewNumberFromFloat(86))
	assert.Equal(t, "undefined: invalid function call - object::*gal_test.Car:SetSpeed - must return 1 value, returned 0 instead", got.String())

	val, ok = gal.ObjectGetMethod(myCar, "SetSpeed2")
	require.True(t, ok)
	got = val(gal.NewNumberFromFloat(96))
	assert.Equal(t, gal.True, got)
	assert.Equal(t, "96", myCar.CurrentSpeed().String())

	val, ok = gal.ObjectGetMethod(myCar, "SetSpeed3")
	require.True(t, ok)
	got = val(gal.NewNumberFromFloat(106))
	assert.Equal(t, gal.True, got)
	assert.Equal(t, "106", myCar.CurrentSpeed().String())

	val, ok = gal.ObjectGetMethod(myCar, "SetSpeed4")
	require.True(t, ok)
	got = val(gal.NewString("blah"))
	assert.Equal(t, "undefined: invalid function call - object::*gal_test.Car:SetSpeed4 - invalid argument type passed to function - reflect: Call using gal.String as type gal_test.fancyType", got.String())
}
