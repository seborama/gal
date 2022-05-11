# Go Eval

A simplistic expression evaluator in Go.

This is a research project.\
It is work in progress and right now in a very early stage.

Check the tests in [`tree_test.go`](tree_test.go) for ideas of usage and capability, notably `TestTree_Eval*` tests.

## Numbers

Numbers implement arbitrary precision fixed-point decimal arithmetic with [shopspring/decimal](https://github.com/shopspring/decimal).

## Supported operations

* Operators: `+` `-` `*` `/` `%` `^`
* Types: String, Number
* Associativity with parentheses
* Functions are syntactically supported but not implemented yet

## High level design

Expressions are parsed in two stages:

- Transformation into a Tree of Values and Operators.
- Evaluation of the Tree for calculation.

Notes:

- a Tree may contain one or more sub-Trees (recursively or not) to hold functions or to express associativity.
- Calculation is performed in successive executions of by decreased operator precedence.
