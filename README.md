# Go Eval

A simplistic expression parser and evaluator in Go.

This is a research project.\
It is work in progress and right now in a very early stage.

## Examples

Check the tests for ideas of usage and capability.

Simple:

```go
func main() {
    expr := `trunc(tan(10 + sin(cos(3*4.4))) 6)`
    gal.Parse(expr).Eval() // returns 3.556049
}
```

With user-defined functions and variables:

```go
// see TestWithVariablesAndFunctions() in gal_test.go for full code
func main() {
	funcs := gal.Functions{
		"double": func(args ...gal.Value) gal.Value {
			// should first validate argument count here
			v := args[0].(gal.Numberer) // should check type assertion is ok here
			return v.Number().Multiply(gal.NewNumber(2))
		},
		"triple": func(args ...gal.Value) gal.Value {
			// should first validate argument count here
			v := args[0].(gal.Numberer)// should check type assertion is ok here
			return v.Number().Multiply(gal.NewNumber(3))
		},
	}

	vars := gal.Variables{
		":val1:": gal.NewNumber(4),
		":val2:": gal.NewNumber(5),
	}

	expr := `double(:val1:) + triple(:val2:)`

	gal.
		Parse(expr, gal.WithFunctions(funcs)).
		Eval(gal.WithVariables(vars)) // returns 23
}
```

## Numbers

Numbers implement arbitrary precision fixed-point decimal arithmetic with [shopspring/decimal](https://github.com/shopspring/decimal).

## Supported operations

* Operators: `+` `-` `*` `/` `%` `^`
* Types: String, Number
* Associativity with parentheses
* Functions:
    * Pre-defined: pi, cos, floor, sin, sqrt, trunc, and more (see `function.go`: `Eval()`)
    * User-defined, injected via `WithFunctions()`
* Variables, defined as `:variable_name:` and injected via `WithVariables()`

## Functions

Function names are case-insensitive.

A function can optionally accept one or more space-separated arguments, but it must return a single Value.

User function definitions are passed as a `map[string]FunctionalValue` using `WithFunctions` when calling `Parse` from `gal`. This may move to `Eval` from `Tree` eventually. This would allow parsing the expression once and run `Eval` multiple times with different function definitions, as with variables.

## Variables

Variable names are case-sensitive.

Values are passed as a `map[string]Value` using `WithVariables` when calling `Eval` from `Tree`.

This allows parsing the expression once with `Parse` and run `Tree`.`Eval` multiple times with different variable values.

## High level design

Expressions are parsed in two stages:

- Transformation into a Tree of Values and Operators.
- Evaluation of the Tree for calculation.

Notes:

- a Tree may contain one or more sub-Trees (recursively or not) to hold functions or to express associativity.
- Calculation is performed in successive rounds of decreased operator precedence. This is to enforce natural associativity.

## Code structure

The main entry point is `Parse` in `gal.go`.

`Parse` instantiates a `TreeBuilder`. It subsequently calls `TreeBuilder`'s `FromExpr` method to create a parsed `Tree` representation of the expression to be evaluated.

Finally, `Tree`'s `Eval` method performs the evaluation of the `Tree` and returns the resultant `Value` to `gal.go`'s `Eval` function.
