# Go Eval

<p align="center">
  <a href="https://pkg.go.dev/github.com/seborama/gal/v8">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg" alt="gal">
  </a>

  <a href="https://goreportcard.com/report/github.com/seborama/gal/v8">
    <img src="https://goreportcard.com/badge/github.com/seborama/gal/v8" alt="gal">
  </a>
</p>

A simple but powerful expression parser and evaluator in Go.

This project started as a personal research.

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

## Type interfaces

`gal` comes with  pre-defined type interfaces: Numberer, Booler, Stringer (and maybe more in the future).

They allow the general use of types. For instance, the String `"123"` can be converted to the Number `123`.
With `Numberer`, a user-defined function can transparently use String and Number when both hold a number representation.

A user-defined function can do this:

```go
n := args[0].(gal.Numberer).Number()
```

or, for additional type safety:

```go
value, ok := args[0].(gal.Numberer)
if !ok {
    return gal.NewUndefinedWithReasonf("NaN '%s'", args[0])
}
n := value.Number()
/* ... */
```

Both examples will happily accept a `Value` of type `String` or `Number` and process it as if it were a `Number`.

## Numbers

Numbers implement arbitrary precision fixed-point decimal arithmetic with [shopspring/decimal](https://github.com/shopspring/decimal).

## Strings

Strings must be enclosed in double-quotes (`"`) e.g. valid: `"this is a string"`, invalid: `this is a syntax error` (missing double-quotes).

Escapes are supported:
- `"this is \"also\" a valid string"`
- `"this is fine too\\"` (escapes cancel each other out)

## Bools

In addition to boolean expressions, sepcial contants `True` and `False` may be used.

Do not double-quote them, or they will become plain strings!

## MultiValue

This is container `Value`. It can contain zero or any number of `Value`'s. Currently, this is only truly useful with functions, mostly because it is yet undecided how to define what operations would mean on a `MultiValue`.

## Supported operations

* Operators: `+` `-` `*` `/` `%` `**` `<<` `>>` `<` `<=` `==` `!=` `>` `>=` `And` `&&` `Or` `||`
    * [Precedence](https://en.wikipedia.org/wiki/Order_of_operations#Programming_languages), highest to lowest:
        * `**`
        * `*` `/` `%`
        * `+` `-`
        * `<<` `>>`
        * `<` `<=` `==` `!=` `>` `>=`
        * `And` `&&` `Or` `||`
    * Notes:
        * Go classifies bit shift operators with the higher `*`.
        * `&&` is synonymous of `And`.
        * `||` is synonymous of `Or`.
        * Worded operators such as `And` and `Or` are **case-sensitive** and must be followed by a blank character. `True Or (False)` is a Bool expression with the `Or` operator but `True Or(False)` is an invalid expression attempting to call a user-defined function called `Or()`.
* Types: String, Number, Bool, MultiValue
* Associativity with parentheses: `(` and `)`
* Functions:
    * Built-in: pi, cos, floor, sin, sqrt, trunc, and more (see `function.go`: `Eval()`)
    * User-defined, injected via `WithFunctions()`
* Variables, defined as `:variable_name:` and injected via `WithVariables()`

## Functions

A function is defined as a Go type: `type FunctionalValue func(...Value) Value`

Function names are case-insensitive.

A **function** can optionally accept one or more **space-separated arguments**, but it must return a single `Value`.

It should be noted that a `MultiValue` type is available that can hold multiple `Value` elements. A function can use `MultiValue` as its return type to effectively return multiple `Value`'s. Of course, as `MultiValue` is a `Value` type, functions can also accept it as part of their argument(s). Refer to the test `TestMultiValueFunctions`, for an example.

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
