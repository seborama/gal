package gal

import "fmt"

type Tree []entry

func (tree Tree) Eval() Value {
	workingTree := tree.CleanUp()

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
			if val == nil && op == invalidOperator {
				val = e.(Value)
				continue
			}

			if val == nil {
				return NewUndefinedWithReason("syntax error: nil value cannot be operated upon (op='" + op.String() + "')")
			}

			rhsVal := e.(Tree).Eval()

			val = calculate(val, op, rhsVal)

		case operatorEntryKind:
			op = e.(Operator)

		case unknownEntryKind:
			// TODO: distinguish between unknownEntryKind and undefinedEntryKind (which is a Value)
			return workingTree[0].(Value)

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
		case treeEntryKind:
			subtree := e.(Tree).cleansePlusMinusTreeStart()
			// TODO: test subtree[0] exists
			if subtree[0].kind() == unknownEntryKind {
				return subtree
			}
			outTree = append(outTree, subtree)
			continue

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

	for i := 0; i < len(tree); i++ {
		e := tree[i]

		switch e.kind() {
		case treeEntryKind:
			subtree := e.(Tree).prioritiseOperators()
			outTree = append(outTree, subtree)
			continue

		case operatorEntryKind:
			outTree = append(outTree, e)

			subTree := associateOnIncreasedPrecedence(tree[i:])
			fmt.Println("DEBUG - i", i, "subTree", subTree)
			if len(subTree) == 0 {
				continue
			}

			outTree = append(outTree, subTree)
			i += len(subTree)

		default:
			outTree = append(outTree, e)
			continue
		}
	}

	return outTree
}

func associateOnIncreasedPrecedence(tree Tree) Tree {
	currentOperatorPrecedence := operatorPrecedence(tree[0].(Operator))

	// fetch the next available operator so to compare precedence and decide
	// on associativity
	nextOperator := invalidOperator

	// TODO: check tree[1] is not out of range
	for _, e := range tree[1:] {
		if e.kind() == operatorEntryKind {
			nextOperator = e.(Operator)
			break
		}
	}

	// when no operator exists in the remainder of the tree (i.e. presumably only a
	// right hand side operand remains) or when the precedence of the next operator is not
	// greater than the current operator, keep processing the tree naturally, left to right,
	// in the current associative context (which may be the root tree).
	if nextOperator == invalidOperator ||
		operatorPrecedence(nextOperator) <= currentOperatorPrecedence {
		return nil
	}

	// the next operator has a greater precedence, start a new sub-tree to creative a new
	// associative context.
	var subTree Tree

	// TODO: check tree[1] is not out of range
	for _, e := range tree[1:] {
		// TODO: if e is a sub-tree, it should be recursively processed by prioritiseOperators too
		if e.kind() == operatorEntryKind &&
			operatorPrecedence(e.(Operator)) != operatorPrecedence(nextOperator) {
			break
		}
		subTree = append(subTree, e)
	}

	return subTree
}

func (Tree) kind() entryKind {
	return treeEntryKind
}
