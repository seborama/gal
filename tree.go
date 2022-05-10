package gal

import (
	"fmt"
)

type Tree []entry

func (tree Tree) TruncLen() int {
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
	fmt.Println("DEBUG start of Eval: tree=", tree)
	workingTree := tree.CleanUp()
	fmt.Println("DEBUG cleaned-up tree:", workingTree)

	var val Value
	var op Operator = invalidOperator

	for i := 0; i < len(workingTree); i++ {
		e := workingTree[i]
		fmt.Println("DEBUG - Eval - i:", i, "e:", e)

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
		panic(fmt.Sprintf("unimplemented operator '%s'", op.String()))
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
	fmt.Println("DEBUG - cleansePlusMinusTreeStart - start - tree:", tree)
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
						outTree = append(outTree, NewNumber(-1), multiply, tree[i+1])
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

	fmt.Println("DEBUG - cleansePlusMinusTreeStart - end - outTree:", outTree)
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

			subTree, lenSubTree := associateOnIncreasedPrecedence(tree[i:])
			fmt.Println("DEBUG - i", i, "subTree:", subTree, "lenSubTree:", lenSubTree)

			i += lenSubTree

			if lenSubTree == 0 {
				continue
			}

			if lenSubTree == 1 { // TODO: is this still necessary?
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

func associateOnIncreasedPrecedence(tree Tree) (Tree, int) {
	fmt.Println("DEBUG start of associateOnIncreasedPrecedence: tree=", tree)
	var outTree Tree
	currentOperator := tree[0].(Operator)

	i := 1
	for ; i < len(tree); i++ {
		nextOperatorIdx := findNextOperator(tree[i:])
		fmt.Println("DEBUG beginning of loop: tree1:", tree)
		fmt.Println("DEBUG beginning of loop: tree1[i:]:", tree[i:])
		fmt.Println("DEBUG beginning of loop: outTree1:", outTree)
		fmt.Println("DEBUG currentOp:", currentOperator, "nextOperatorIdx=", nextOperatorIdx)

		if nextOperatorIdx == noOpFound {
			outTree = append(outTree, tree[i:]...)
			break
		}

		nextOperatorIdx += i
		nextOperator := tree[nextOperatorIdx].(Operator)

		pComp := compareOperatorPrecedence(currentOperator, nextOperator)
		if pComp < 0 {
			// next operator is of lesser or equal precedence
			outTree = append(outTree, tree[i:nextOperatorIdx]...)
			fmt.Println("DEBUG currentOp:", currentOperator, "pComp:", pComp, "outTree=", outTree)
			fmt.Println("DEBUG currentOp:", currentOperator, "pComp:", pComp, "tree[i:nextOperatorIdx]=", tree[i:nextOperatorIdx])
			break
		}

		if pComp == 0 {
			// next operator is of equal precedence
			outTree = append(outTree, tree[i:nextOperatorIdx+1]...)
			fmt.Println("DEBUG currentOp:", currentOperator, "pComp:", pComp, "outTree=", outTree)
			fmt.Println("DEBUG currentOp:", currentOperator, "pComp:", pComp, "tree[i:nextOperatorIdx+1]=", tree[i:nextOperatorIdx+1])
			continue
		}

		subTree := Tree{}
		subTree = append(subTree, tree[i:nextOperatorIdx+1]...)
		i += len(subTree)
		fmt.Println("DEBUG currentOp:", currentOperator, "subTree1=", subTree, "nextOp", nextOperator)

		moreTreeParts, lenMoreTreeParts := associateOnIncreasedPrecedence(tree[nextOperatorIdx:])
		fmt.Println("DEBUG outTree2:", outTree)
		fmt.Println("DEBUG moreTreeParts:", moreTreeParts, "lenMoreTreeParts:", lenMoreTreeParts)
		fmt.Println("DEBUG currentOp:", currentOperator, "subTree2=", subTree, "nextOp", nextOperator)

		subTree = append(subTree, moreTreeParts)
		fmt.Println("DEBUG currentOp:", currentOperator, "subTree3=", subTree, "nextOp", nextOperator)

		outTree = append(outTree, subTree)
		fmt.Println("DEBUG outTree3:", outTree, "subTree.TruncLen():", subTree.TruncLen(), "subTree.FullLen():", subTree.FullLen())
		i += lenMoreTreeParts - 1
	}

	fmt.Println("DEBUG outTree (end):", outTree)
	return outTree, i
}

func associateOnIncreasedPrecedenceV1(tree Tree) Tree {
	fmt.Println("DEBUG start of associateOnIncreasedPrecedence: tree=", tree)
	currentOperator := tree[0].(Operator)

	shift := 1
	nextOperatorIdx := findNextOperator(tree[shift:]) // TODO: check tree[1] is not out of range

	if nextOperatorIdx == -1 {
		// no operator in the remainder of the tree
		resTree := Tree{}
		resTree = append(resTree, tree[shift:]...)
		return resTree
	}

	nextOperatorIdx += shift
	nextOperator := tree[nextOperatorIdx].(Operator)

	pComp := compareOperatorPrecedence(currentOperator, nextOperator)
	if pComp <= 0 {
		// next operator is of lesser or equal precedence
		resTree := Tree{}
		resTree = append(resTree, tree[shift:nextOperatorIdx]...)
		return resTree
	}

	subTree := Tree{}
	subTree = append(subTree, tree[shift:nextOperatorIdx+1]...)
	fmt.Println("DEBUG currentOp:", currentOperator, "subTree1=", subTree, "nextOp", nextOperator)

	moreTreeParts := associateOnIncreasedPrecedenceV1(tree[nextOperatorIdx:])
	fmt.Println("DEBUG currentOp:", currentOperator, "subTree2=", subTree, "nextOp", nextOperator)

	subTree = append(subTree, moreTreeParts)
	fmt.Println("DEBUG currentOp:", currentOperator, "subTree3=", subTree, "nextOp", nextOperator)

	return subTree
}

// returns -1 if none found
const noOpFound = -1

func findNextOperator(tree Tree) int {
	pos := 0

	for _, e := range tree[pos:] {
		if e.kind() == operatorEntryKind {
			return pos
		}
		pos++
	}

	return noOpFound
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
