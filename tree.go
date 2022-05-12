package gal

import (
	"fmt"
)

type Tree []entry

func (tree Tree) TrunkLen() int {
	return len(tree)
}

func (tree Tree) FullLen() int {
	l := len(tree)

	for _, e := range tree {
		if subTree, ok := e.(Tree); ok {
			l += subTree.FullLen() - 1
		}
	}

	return l
}

func (tree Tree) Eval() Value {
	// execute calculation by decreasing order of precedence
	workingTree := tree.CleanUp().
		Calc(powerOperators).
		Calc(multiplicativeOperators).
		Calc(additiveOperators)

	// TODO: refactor this
	// perhaps add Tree.Value() which tests that only one entry is left and that it is a Value
	return workingTree[0].(Value)
}

func (tree Tree) Calc(isOperatorInFocus func(Operator) bool) Tree {
	var outTree Tree

	var val entry
	var op Operator = invalidOperator

	for i := 0; i < tree.TrunkLen(); i++ {
		e := tree[i]

		switch e.kind() {
		case valueEntryKind:
			if val == nil && op == invalidOperator {
				val = e
				continue
			}

			if val == nil {
				return Tree{
					NewUndefinedWithReason("syntax error: nil value cannot be operated upon (op='" + op.String() + "')"),
				}
			}

			val = calculate(val.(Value), op, e.(Value))

		case treeEntryKind:
			if val == nil && op != invalidOperator {
				return Tree{
					NewUndefinedWithReason("syntax error: nil value cannot be operated upon (op='" + op.String() + "')"),
				}
			}

			rhsVal := e.(Tree).Eval()
			if val == nil {
				val = rhsVal
				continue
			}

			val = calculate(val.(Value), op, rhsVal)

		case operatorEntryKind:
			op = e.(Operator)
			if isOperatorInFocus(op) {
				continue
			}
			outTree = append(outTree, val)
			outTree = append(outTree, op)
			val = nil
			op = invalidOperator

		case unknownEntryKind:
			// TODO: distinguish between unknownEntryKind and undefinedEntryKind (which is a Value)
			return Tree{e}

		default:
			return Tree{
				NewUndefinedWithReason(fmt.Sprintf("internal error: unknown entry kind: '%v'", e.kind())),
			}
		}
	}

	if val != nil {
		outTree = append(outTree, val)
	}

	return outTree
}

func calculate(lhs Value, op Operator, rhs Value) Value {
	var outVal Value = NewUndefined()

	switch op {
	case plus:
		outVal = lhs.Add(rhs)

	case minus:
		outVal = lhs.Sub(rhs)

	case multiply:
		outVal = lhs.Multiply(rhs)

	case divide:
		outVal = lhs.Divide(rhs)

	case power:
		outVal = lhs.PowerOf(rhs)

	case modulus:
		outVal = lhs.Mod(rhs)

	default:
		return NewUndefinedWithReason(fmt.Sprintf("unimplemented operator: '%s'", op.String()))
	}

	return outVal
}

func (tree Tree) CleanUp() Tree {
	// TODO: add syntaxCheck() before prioritiseOperators to obtain a syntactically correct tree
	//       i.e. 3 * * 4 would be detected as a syntax error, etc
	return tree.
		cleansePlusMinusTreeStart()
}

// cleansePlusMinusTreeStart consolidates the - and + that are at the first position in a Tree.
// `plus` is removed and `minus` causes the number that follows to be negated.
func (tree Tree) cleansePlusMinusTreeStart() Tree {
	outTree := Tree{}

	for i := 0; i < len(tree); i++ {
		e := tree[i]

		switch e.kind() {
		case operatorEntryKind:
			if i != 0 {
				outTree = append(outTree, e)
				continue
			}

			switch e.(Operator) {
			case plus:
				// drop superfluous plus sign at start of Tree
				continue

			case minus:
				outTree = append(outTree, NewNumber(-1), multiply, tree[i+1])
				i++
				continue

			default:
				return Tree{
					NewUndefinedWithReason("syntax error: expression starts with '" + e.(Operator).String() + "'"),
				}
			}

		default:
			outTree = append(outTree, e)
		}
	}

	return outTree
}

func (Tree) kind() entryKind {
	return treeEntryKind
}
