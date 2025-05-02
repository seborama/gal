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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
