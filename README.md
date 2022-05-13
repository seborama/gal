# Go Eval

A simplistic expression evaluator in Go.

This is a research project.\
It is work in progress and right now in a very early stage.

Check the tests for ideas of usage and capability.

## Numbers

Numbers implement arbitrary precision fixed-point decimal arithmetic with [shopspring/decimal](https://github.com/shopspring/decimal).

## Supported operations

* Operators: `+` `-` `*` `/` `%` `^`
* Types: String, Number
* Associativity with parentheses
* Functions:
    * Pre-defined: pi, cos, floor, sin, sqrt, trunc, and more (see `function.go`: `Eval()`)
    * User-defined: TODO
* Variables, defined as `:variable_name:`

## Functions

Function names are case-insensitive.
A function can optionally accept one or more arguments but it must return a single Value.

## Variables

Variable names are case-sensitive.

Values are passed as a `map[string]Value` using `WithVariables` when calling `Eval`.

## High level design

Expressions are parsed in two stages:

- Transformation into a Tree of Values and Operators.
- Evaluation of the Tree for calculation.

Notes:

- a Tree may contain one or more sub-Trees (recursively or not) to hold functions or to express associativity.
- Calculation is performed in successive rounds of decreased operator precedence. This is to enforce natural associativity.

## Code structure

The main entry point is `Eval` in `gal.go`.

`Eval` instantiates a `TreeBuilder` optionally with configuration (notably to pass a map of variable names and values). It subsequently calls `TreeBuilder`'s `FromExpr` method to create a `Tree` from the expression to be evaluated.

Finally, `Tree`'s `Eval` method performs the evaluation of the `Tree` and returns the resultant `Value` to `gal.go`'s `Eval` function.
