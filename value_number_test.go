package gal

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewNumber(t *testing.T) {
	got := NewNumber(5, -2)
	assert.Equal(t, Number{Undefined: Undefined{}, value: decimal.New(5, -2)}, got)
}

func TestNewNumberFromInt(t *testing.T) {
	got := NewNumberFromInt(5)
	assert.Equal(t, Number{Undefined: Undefined{}, value: decimal.New(5, 0)}, got)
}

func TestNewNumberFromFloat(t *testing.T) {
	got := NewNumberFromFloat(5.45678)
	assert.Equal(t, Number{Undefined: Undefined{}, value: decimal.NewFromFloat(5.45678)}, got)
}

func TestNewNumberFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Number
		wantErr bool
	}{
		{
			name:    "it creates a number from a string",
			args:    args{s: "5.45678"},
			want:    Number{Undefined: Undefined{}, value: decimal.NewFromFloat(5.45678)},
			wantErr: false,
		},
		{
			name:    "it returns an error when the string is not a number",
			args:    args{s: "not a number"},
			want:    Number{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNumberFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNumberFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNumberFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Equal(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Number
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "equal",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(5, 0)}},
			want:   true,
		},
		{
			name:   "not equal",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(6, 0)}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Equal(tt.args.other); got != tt.want {
				t.Errorf("Number.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Add(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		{
			name:   "add two positive numbers",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(3, 0)}},
			want:   Number{value: decimal.New(8, 0)},
		},
		{
			name:   "add a positive and a negative number",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(2, 0)},
		},
		{
			name:   "add two negative numbers",
			fields: fields{value: decimal.New(-5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(-8, 0)},
		},
		{
			name:   "add a number and non-Numberer",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: NewMultiValue()},
			want:   NewUndefinedWithReasonf("NaN: %s", MultiValue{}.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Add(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Sub(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		{
			name:   "subtract two positive numbers",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(3, 0)}},
			want:   Number{value: decimal.New(2, 0)},
		},
		{
			name:   "subtract a positive and a negative number",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(8, 0)},
		},
		{
			name:   "subtract two negative numbers",
			fields: fields{value: decimal.New(-5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(-2, 0)},
		},
		{
			name:   "subtract a number and non-Numberer",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: NewMultiValue()},
			want:   NewUndefinedWithReasonf("NaN: %s", MultiValue{}.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Sub(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Multiply(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		{
			name:   "multiply two positive numbers",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(3, 0)}},
			want:   Number{value: decimal.New(15, 0)},
		},
		{
			name:   "multiply a positive and a negative number",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(-15, 0)},
		},
		{
			name:   "multiply two negative numbers",
			fields: fields{value: decimal.New(-5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(15, 0)},
		},
		{
			name:   "multiply a number and non-Numberer",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: NewMultiValue()},
			want:   NewUndefinedWithReasonf("NaN: %s", MultiValue{}.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Multiply(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Multiply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Divide(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		{
			name:   "divide two positive numbers",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(3, 0)}},
			want:   Number{value: decimal.New(5, 0).Div(decimal.New(3, 0))},
		},
		{
			name:   "divide a positive and a negative number",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(5, 0).Div(decimal.New(-3, 0))},
		},
		{
			name:   "divide two negative numbers",
			fields: fields{value: decimal.New(-5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(-5, 0).Div(decimal.New(-3, 0))},
		},
		{
			name:   "divide a number and non-Numberer",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: NewMultiValue()},
			want:   NewUndefinedWithReasonf("NaN: %s", MultiValue{}.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Divide(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Divide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_PowerOf(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		{
			name:   "power of two positive numbers",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(3, 0)}},
			want:   Number{value: decimal.New(125, 0)},
		},
		{
			name:   "power of a positive and a negative number",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(8e13, -16)},
		},
		{
			name:   "power of two negative numbers",
			fields: fields{value: decimal.New(-5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(-8e13, -16)},
		},
		{
			name:   "power of a number and non-Numberer",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: NewMultiValue()},
			want:   NewUndefinedWithReasonf("NaN: %s", MultiValue{}.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.PowerOf(tt.args.other); !assert.Equal(t, tt.want, got) {
				t.Errorf("Number.PowerOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Mod(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		{
			name:   "modulus of two positive numbers",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(3, 0)}},
			want:   Number{value: decimal.New(2, 0)},
		},
		{
			name:   "modulus of a positive and a negative number",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(2, 0)},
		},
		{
			name:   "modulus of two negative numbers",
			fields: fields{value: decimal.New(-5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Number{value: decimal.New(-2, 0)},
		},
		{
			name:   "modulus of a number and non-Numberer",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: NewMultiValue()},
			want:   NewUndefinedWithReasonf("NaN: %s", MultiValue{}.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Mod(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_IntPart(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   Value
	}{
		{
			name:   "integer part of a positive number",
			fields: fields{value: decimal.New(5678, -3)}, // 5.678
			want:   Number{value: decimal.New(5, 0)},
		},
		{
			name:   "integer part of a negative number",
			fields: fields{value: decimal.New(-5678, -3)}, // -5.678
			want:   Number{value: decimal.New(-5, 0)},
		},
		{
			name:   "integer part of zero",
			fields: fields{value: decimal.New(0, 0)},
			want:   Number{value: decimal.New(0, 0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.IntPart(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.IntPart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_LShift(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		{
			name:   "left shift two positive numbers",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(3, 0)}},
			want:   Number{value: decimal.New(40, 0)},
		},
		{
			name:   "left shift a negative and a positive number",
			fields: fields{value: decimal.New(-5, 0)},
			args:   args{other: Number{value: decimal.New(3, 0)}},
			want:   Number{value: decimal.New(-40, 0)},
		},
		{
			name:   "left shift a positive and a negative number",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Undefined{"invalid negative left shift"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.LShift(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.LShift() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_RShift(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		{
			name:   "right shift two positive numbers",
			fields: fields{value: decimal.New(500, 0)},
			args:   args{other: Number{value: decimal.New(3, 0)}},
			want:   Number{value: decimal.New(62, 0)},
		},
		{
			name:   "right shift a negative and a positive number",
			fields: fields{value: decimal.New(-500, 0)},
			args:   args{other: Number{value: decimal.New(3, 0)}},
			want:   Number{value: decimal.New(-63, 0)},
		},
		{
			name:   "left shift a positive and a negative number",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(-3, 0)}},
			want:   Undefined{"invalid negative right shift"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.RShift(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.RShift() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Neg(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   Number
	}{
		{
			name:   "negate a positive number",
			fields: fields{value: decimal.New(5, 0)},
			want:   Number{value: decimal.New(-5, 0)},
		},
		{
			name:   "negate a negative number",
			fields: fields{value: decimal.New(-5, 0)},
			want:   Number{value: decimal.New(5, 0)},
		},
		{
			name:   "negate zero",
			fields: fields{value: decimal.New(0, 0)},
			want:   Number{value: decimal.New(0, 0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Neg(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Neg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Sin(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   Number
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Sin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Sin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Cos(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   Number
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Cos(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Cos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Sqrt(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   Value
	}{
		{
			name:   "sqrt of a positive number",
			fields: fields{value: decimal.New(927344, 0)},
			want:   Number{value: decimal.NewFromFloat(962.9870196)},
		},
		{
			name:   "sqrt of zero",
			fields: fields{value: decimal.New(0, 0)},
			want:   Number{value: decimal.New(0, 0)},
		},
		{
			name:   "sqrt of a negative number",
			fields: fields{value: decimal.New(-4, 0)},
			want:   Undefined{"square root of negative number: -4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Sqrt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Sqrt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Tan(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   Number
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Tan(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Tan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Ln(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		precision int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Ln(tt.args.precision); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Ln() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Log(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		precision int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Log(tt.args.precision); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Log() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Floor(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   Number
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Floor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Floor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Trunc(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		precision int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Number
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Trunc(tt.args.precision); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Trunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Factorial(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   Value
	}{
		{
			name:   "factorial of a positive number",
			fields: fields{value: decimal.New(5, 0)},
			want:   Number{value: decimal.New(120, 0)},
		},
		{
			name:   "factorial of zero",
			fields: fields{value: decimal.New(0, 0)},
			want:   Number{value: decimal.New(1, 0)},
		},
		{
			name:   "factorial of a negative number",
			fields: fields{value: decimal.New(-5, 0)},
			want:   Undefined{"Factorial: requires a positive integer, cannot accept -5"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Factorial(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Factorial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_LessThan(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Bool
	}{
		{
			name:   "less than",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(6, 0)}},
			want:   Bool{value: true},
		},
		{
			name:   "not less than",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(4, 0)}},
			want:   Bool{value: false},
		},
		{
			name:   "equal",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(5, 0)}},
			want:   Bool{value: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.LessThan(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.LessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_LessThanOrEqual(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Bool
	}{
		{
			name:   "less than or equal",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(6, 0)}},
			want:   Bool{value: true},
		},
		{
			name:   "not less than or equal",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(4, 0)}},
			want:   Bool{value: false},
		},
		{
			name:   "equal",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(5, 0)}},
			want:   Bool{value: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.LessThanOrEqual(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.LessThanOrEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_EqualTo(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Bool
	}{
		{
			name:   "equal to",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(5, 0)}},
			want:   Bool{value: true},
		},
		{
			name:   "not equal to",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(4, 0)}},
			want:   Bool{value: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.EqualTo(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.EqualTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_NotEqualTo(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Bool
	}{
		{
			name:   "not equal to",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(4, 0)}},
			want:   Bool{value: true},
		},
		{
			name:   "equal to",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(5, 0)}},
			want:   Bool{value: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.NotEqualTo(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.NotEqualTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_GreaterThan(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Bool
	}{
		{
			name:   "greater than",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(4, 0)}},
			want:   Bool{value: true},
		},
		{
			name:   "not greater than",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(6, 0)}},
			want:   Bool{value: false},
		},
		{
			name:   "equal",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(5, 0)}},
			want:   Bool{value: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.GreaterThan(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.GreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_GreaterThanOrEqual(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	type args struct {
		other Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Bool
	}{
		{
			name:   "greater than or equal",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(4, 0)}},
			want:   Bool{value: true},
		},
		{
			name:   "not greater than or equal",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(6, 0)}},
			want:   Bool{value: false},
		},
		{
			name:   "equal",
			fields: fields{value: decimal.New(5, 0)},
			args:   args{other: Number{value: decimal.New(5, 0)}},
			want:   Bool{value: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.GreaterThanOrEqual(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.GreaterThanOrEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_String(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "string representation of a positive number",
			fields: fields{value: decimal.New(5, 0)},
			want:   "5",
		},
		{
			name:   "string representation of a negative number",
			fields: fields{value: decimal.New(-5, 0)},
			want:   "-5",
		},
		{
			name:   "string representation of zero",
			fields: fields{value: decimal.New(0, 0)},
			want:   "0",
		},
		{
			name:   "string representation of a decimal number",
			fields: fields{value: decimal.New(5, -2)},
			want:   "0.05",
		},
		{
			name:   "string representation of a large number",
			fields: fields{value: decimal.New(1234567890123456789, 0)},
			want:   "1234567890123456789",
		},
		{
			name:   "string representation of a small number",
			fields: fields{value: decimal.New(1234567890123456789, -20)},
			want:   "0.01234567890123456789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.String(); got != tt.want {
				t.Errorf("Number.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Bool(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   Bool
	}{
		{
			name:   "true for non-zero number",
			fields: fields{value: decimal.New(5, 0)},
			want:   Bool{value: true},
		},
		{
			name:   "false for zero",
			fields: fields{value: decimal.New(0, 0)},
			want:   Bool{value: false},
		},
		{
			name:   "true for negative number",
			fields: fields{value: decimal.New(-5, 0)},
			want:   Bool{value: true},
		},
		{
			name:   "true for positive decimal number",
			fields: fields{value: decimal.New(5, -2)},
			want:   Bool{value: true},
		},
		{
			name:   "false for zero decimal number",
			fields: fields{value: decimal.New(0, -2)},
			want:   Bool{value: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Bool(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_AsString(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   String
	}{
		{
			name:   "string representation of a positive number",
			fields: fields{value: decimal.New(5, 0)},
			want:   String{value: "5"},
		},
		{
			name:   "string representation of a negative number",
			fields: fields{value: decimal.New(-5, 0)},
			want:   String{value: "-5"},
		},
		{
			name:   "string representation of zero",
			fields: fields{value: decimal.New(0, 0)},
			want:   String{value: "0"},
		},
		{
			name:   "string representation of a decimal number",
			fields: fields{value: decimal.New(5, -2)},
			want:   String{value: "0.05"},
		},
		{
			name:   "string representation of a large number",
			fields: fields{value: decimal.New(1234567890123456789, 0)},
			want:   String{value: "1234567890123456789"},
		},
		{
			name:   "string representation of a small number",
			fields: fields{value: decimal.New(1234567890123456789, -20)},
			want:   String{value: "0.01234567890123456789"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.AsString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Number(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   Number
	}{
		{
			name:   "number representation of a positive number",
			fields: fields{value: decimal.New(5, 0)},
			want:   Number{value: decimal.New(5, 0)},
		},
		{
			name:   "number representation of a negative number",
			fields: fields{value: decimal.New(-5, 0)},
			want:   Number{value: decimal.New(-5, 0)},
		},
		{
			name:   "number representation of zero",
			fields: fields{value: decimal.New(0, 0)},
			want:   Number{value: decimal.New(0, 0)},
		},
		{
			name:   "number representation of a decimal number",
			fields: fields{value: decimal.New(5, -2)},
			want:   Number{value: decimal.New(5, -2)},
		},
		{
			name:   "number representation of a large number",
			fields: fields{value: decimal.New(1234567890123456789, 0)},
			want:   Number{value: decimal.New(1234567890123456789, 0)},
		},
		{
			name:   "number representation of a small number",
			fields: fields{value: decimal.New(1234567890123456789, -20)},
			want:   Number{value: decimal.New(1234567890123456789, -20)},
		},
		{
			name:   "number representation of a large decimal number",
			fields: fields{value: decimal.New(1234567890123456789, -10)},
			want:   Number{value: decimal.New(1234567890123456789, -10)},
		},
		{
			name:   "number representation of a small decimal number",
			fields: fields{value: decimal.New(1234567890123456789, -30)},
			want:   Number{value: decimal.New(1234567890123456789, -30)},
		},
		{
			name:   "number representation of a large negative number",
			fields: fields{value: decimal.New(-1234567890123456789, 0)},
			want:   Number{value: decimal.New(-1234567890123456789, 0)},
		},
		{
			name:   "number representation of a small negative number",
			fields: fields{value: decimal.New(-1234567890123456789, -20)},
			want:   Number{value: decimal.New(-1234567890123456789, -20)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Number(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Number() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Float64(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name:   "float64 representation of a positive number",
			fields: fields{value: decimal.New(5, 0)},
			want:   5.0,
		},
		{
			name:   "float64 representation of a negative number",
			fields: fields{value: decimal.New(-5, 0)},
			want:   -5.0,
		},
		{
			name:   "float64 representation of zero",
			fields: fields{value: decimal.New(0, 0)},
			want:   0.0,
		},
		{
			name:   "float64 representation of a decimal number",
			fields: fields{value: decimal.New(5, -2)},
			want:   0.05,
		},
		{
			name:   "float64 representation of a large number",
			fields: fields{value: decimal.New(1234567890123456789, 0)},
			want:   1234567890123456789.0,
		},
		{
			name:   "float64 representation of a small number",
			fields: fields{value: decimal.New(1234567890123456789, -20)},
			want:   0.01234567890123456789,
		},
		{
			name:   "float64 representation of a large decimal number",
			fields: fields{value: decimal.New(1234567890123456789, -10)},
			want:   1.2345678901234567e+08, // rounded, may be architecture dependent
		},
		{
			name:   "float64 representation of a small decimal number",
			fields: fields{value: decimal.New(1234567890123456789, -30)},
			want:   1.2345678901234569e-12, // rounded, may be architecture dependent
		},
		{
			name:   "float64 representation of a large negative number",
			fields: fields{value: decimal.New(-1234567890123456789, 0)},
			want:   -1234567890123456789.0,
		},
		{
			name:   "float64 representation of a small negative number",
			fields: fields{value: decimal.New(-1234567890123456789, -20)},
			want:   -0.01234567890123456789,
		},
		{
			name:   "float64 representation of a large negative decimal number",
			fields: fields{value: decimal.New(-1234567890123456789, -10)},
			want:   -1.2345678901234567e+08, // rounded, may be architecture dependent
		},
		{
			name:   "float64 representation of a small negative decimal number",
			fields: fields{value: decimal.New(-1234567890123456789, -30)},
			want:   -1.2345678901234569e-12, // rounded, may be architecture dependent
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Float64(); got != tt.want {
				t.Errorf("Number.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Int64(t *testing.T) {
	type fields struct {
		Undefined Undefined
		value     decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name:   "int64 representation of a positive number",
			fields: fields{value: decimal.New(5, 0)},
			want:   5,
		},
		{
			name:   "int64 representation of a negative number",
			fields: fields{value: decimal.New(-5, 0)},
			want:   -5,
		},
		{
			name:   "int64 representation of zero",
			fields: fields{value: decimal.New(0, 0)},
			want:   0,
		},
		{
			name:   "int64 representation of a decimal number",
			fields: fields{value: decimal.New(5, -2)},
			want:   0,
		},
		{
			name:   "int64 representation of a large number",
			fields: fields{value: decimal.New(1234567890123456789, 0)},
			want:   1234567890123456789,
		},
		{
			name:   "int64 representation of a small number",
			fields: fields{value: decimal.New(1234567890123456789, -20)},
			want:   0,
		},
		{
			name:   "int64 representation of a large decimal number",
			fields: fields{value: decimal.New(1234567890123456789, -10)},
			want:   123456789,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number{
				Undefined: tt.fields.Undefined,
				value:     tt.fields.value,
			}
			if got := n.Int64(); got != tt.want {
				t.Errorf("Number.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}
