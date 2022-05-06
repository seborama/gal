package gal

import (
	"fmt"
)

type Tree []entry

func (tree Tree) Eval() Value {
	fmt.Println("DEBUG start of Eval: tree=", tree)
	workingTree := tree.CleanUp()
	fmt.Println("DEBUG cleaned-up tree:", tree)

	var val Value
	var op Operator = invalidOperator

	for i := 0; i < len(workingTree); i++ {
		e := workingTree[i]

		switch e.kind() {
		case valueEntryKind:
			if val == nil && op == invalidOperator {
				val = e.(Value)
				continue
			}

			if val == nil {
				return NewUndefinedWithReason("syntax error: nil value cannot be operated upon (op='" + op.String() + "')")
			}

			val = calculate(val, op, e.(Value))

		case treeEntryKind:
			if val == nil && op != invalidOperator {
				return NewUndefinedWithReason("syntax error: nil value cannot be operated upon (op='" + op.String() + "')")
			}

			rhsVal := e.(Tree).Eval()
			if val == nil {
				val = rhsVal
				continue
			}

			val = calculate(val, op, rhsVal)

		case operatorEntryKind:
			op = e.(Operator)

		case unknownEntryKind:
			// TODO: distinguish between unknownEntryKind and undefinedEntryKind (which is a Value)
			return e.(Value)

		default:
			panic("TODO")
		}
	}

	return val
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

	case power:
		outVal = lhs.PowerOf(rhs)

	default:
		panic("TODO")
	}

	return outVal
}

func (tree Tree) CleanUp() Tree {
	// TODO: add syntaxCheck() before prioritiseOperators to obtain a syntactically correct tree
	//       i.e. 3 * * 4 would be detected as a syntax error, etc
	return tree.
		cleansePlusMinusTreeStart().
		prioritiseOperators()
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
				// TODO: check that i+1 is not out of range
				if tree[i+1].kind() == valueEntryKind {
					if _, ok := tree[i+1].(Number); ok {
						newEntry := tree[i+1].(Number).Neg()
						outTree = append(outTree, newEntry)
						i++
						continue
					}

					// TODO: is checking for all these syntax error conditions here truly necessary?
					return Tree{
						NewUndefinedWithReason("syntax error: attempt to negate non-number at start of tree"),
					}
				}

				return Tree{
					NewUndefinedWithReason("syntax error: attempt to negate non-value at start of tree"),
				}

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

// prioritiseOperators walks the tree and moves portions to a sub-tree to enforce operator
// precedence through associativity.
// Note: Tree is expected to be syntactically safe. Support for expressions such as `1+2+3 4*5`
//       should be implemented by the Tree semantics, not by this function alone. For instance,
//       this could be through splitting expressions that produce multiple values as []Tree.
func (tree Tree) prioritiseOperators() Tree {
	outTree := Tree{}
	fmt.Println("DEBUG start of prioritiseOperators: tree=", tree)

	for i := 0; i < len(tree); i++ {
		e := tree[i]
		fmt.Println("DEBUG - i=", i, "e=", e, "outTree=", outTree)

		switch e.kind() {
		case operatorEntryKind:
			outTree = append(outTree, e)

			subTree := associateOnIncreasedPrecedence(tree[i:])
			fmt.Println("DEBUG - i", i, "subTree", subTree)

			lenSubTree := len(subTree)
			i += len(subTree)

			if lenSubTree == 0 {
				continue
			}

			if lenSubTree == 1 { // TODO: is thiis still necessary?
				outTree = append(outTree, subTree[0])
				continue
			}

			outTree = append(outTree, subTree)

		default:
			outTree = append(outTree, e)
			continue
		}
	}

	fmt.Println("DEBUG end of prioritiseOperators: outTree=", outTree)
	return outTree
}

func associateOnIncreasedPrecedence(tree Tree) Tree {
	fmt.Println("DEBUG start of associateOnIncreasedPrecedence: tree=", tree)
	currentOperator := tree[0].(Operator)

	shift := 1
	nextOperatorIdx := findNextOperator(tree[shift:]) // TODO: check tree[1] is not out of range

	if nextOperatorIdx == -1 {
		// no operator in the remainder of the tree
		return tree[shift:]
	}

	nextOperatorIdx += shift
	nextOperator := tree[nextOperatorIdx].(Operator)

	pComp := compareOperatorPrecedence(currentOperator, nextOperator)
	if pComp <= 0 {
		// next operator is of lesser or equal precedence
		return tree[shift:nextOperatorIdx]
	}

	subTree := tree[shift : nextOperatorIdx+1]
	fmt.Println("DEBUG cuurenOp:", currentOperator, "subTree1=", subTree, "nextOp", nextOperator)

	moreTreeParts := associateOnIncreasedPrecedence(tree[nextOperatorIdx:])
	subTree = append(subTree, moreTreeParts...)
	fmt.Println("DEBUG subTree2=", subTree)

	return subTree
}

// returns -1 if none found
func findNextOperator(tree Tree) int {
	pos := 0

	for _, e := range tree[pos:] {
		if e.kind() == operatorEntryKind {
			return pos
		}
		pos++
	}

	return -1
}

// returns -1 if operator b has a lesser precedence than operator a.
// returns 0 if both operators have the same precedence.
// returns 1 if operator b has a greater precedence than operator a.
func compareOperatorPrecedence(a, b Operator) int {
	precA := operatorPrecedence(a)
	precB := operatorPrecedence(b)

	if precA == precB {
		return 0
	}

	if precB < precA {
		return -1
	}

	return 1
}

func (Tree) kind() entryKind {
	return treeEntryKind
}
