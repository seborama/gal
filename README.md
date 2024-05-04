# Go Eval

<p align="center">
  <a href="https://pkg.go.dev/github.com/seborama/gal/v7">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg" alt="gal">
  </a>

  <a href="https://goreportcard.com/report/github.com/seborama/gal/v7">
    <img src="https://goreportcard.com/badge/github.com/seborama/gal/v7" alt="gal">
  </a>
</p>

A simple but powerful expression parser and evaluator in Go.

This is a research project.

Short link to repo Readme:

## Examples

Check the tests for ideas of usage and capability.

Simple:

```go
func main() {
    expr := `trunc(tan(10 + sin(cos(3*4.4))) 6)`
    gal.Parse(expr).Eval() // returns 3.556049
}
```

Advanced example, with user-defined functions and variables redefined once.\
In this case, the expression is parsed once but evaluate twice:

```go
// see TestWithVariablesAndFunctions() in gal_test.go for full code
func main() {
    // first of all, parse the expression (once only)
    expr := `double(:val1:) + triple(:val2:)`
    parsedExpr := gal.Parse(expr)

    // step 1: define funcs and vars and Eval the expression
    funcs := gal.Functions{
        "double": func(args ...gal.Value) gal.Value {
            value := args[0].(gal.Numberer)
            return value.Number().Multiply(gal.NewNumber(2))
        },
        "triple": func(args ...gal.Value) gal.Value {
            value := args[0].(gal.Numberer)
            return value.Number().Multiply(gal.NewNumber(3))
        },
    }

    vars := gal.Variables{
        ":val1:": gal.NewNumber(4),
        ":val2:": gal.NewNumber(5),
    }

    // returns 4 * 2 + 5 * 3 == 23
    parsedExpr.Eval(
        gal.WithVariables(vars),
        gal.WithFunctions(funcs),
    )

    // step 2: re-define funcs and vars and Eval the expression again
    // note that we do not need to parse the expression again, only just evaluate it
    funcs = gal.Functions{
        "double": func(args ...gal.Value) gal.Value {
            value := args[0].(gal.Numberer)
            return value.Number().Divide(gal.NewNumber(2))
        },
        "triple": func(args ...gal.Value) gal.Value {
            value := args[0].(gal.Numberer)
            return value.Number().Divide(gal.NewNumber(3))
        },
    }

    vars = gal.Variables{
        ":val1:": gal.NewNumber(2),
        ":val2:": gal.NewNumber(6),
    }

    // returns 2 / 2 + 6 / 3 == 3 this time
    parsedExpr.Eval(
        gal.WithVariables(vars),
        gal.WithFunctions(funcs),
    )
}
```

## Numbers

Numbers implement arbitrary precision fixed-point decimal arithmetic with [shopspring/decimal](https://github.com/shopspring/decimal).

## Supported operations

* Operators: `+` `-` `*` `/` `%` `**` `<<` `>>` `<` `<=` `==` `!=` `>` `>=`
    * [Precedence](https://en.wikipedia.org/wiki/Order_of_operations#Programming_languages), highest to lowest:
        * `**`
        * `*` `/` `%`
        * `+` `-`
        * `<<` `>>`
        * `<` `<=` `==` `!=` `>` `>=`
    * Note: Go classifies bit shift operators with the higher `*`.
* Types: String, Number, Bool
* Associativity with parentheses
* Functions:
    * Pre-defined: pi, cos, floor, sin, sqrt, trunc, and more (see `function.go`: `Eval()`)
    * User-defined, injected via `WithFunctions()`
* Variables, defined as `:variable_name:` and injected via `WithVariables()`

## Functions

Function names are case-insensitive.

A function can optionally accept one or more space-separated arguments, but it must return a single Value.

User function definitions are passed as a `map[string]FunctionalValue` using `WithFunctions` when calling `Eval` from `Tree`.

This allows parsing the expression once with `Parse` and run `Tree`.`Eval` multiple times with different user function definitions.

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

## To do

A number of TODO's exist throughout the code.

The next priorities are:
- review `TODO`'s
