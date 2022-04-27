package gal

type Tree []entry

func (tree Tree) Eval() Value {
	for i := 0; i < len(tree); i++ {
		e := tree[i]

		switch e.kind() {
		case treeEntryKind:
		case operatorEntryKind:
		default:
		}
	}
	panic("TODO")
}

func (tree Tree) PrioritiseOperators() Tree {
	// TODO perhaps more functions like this one needed to deal with leading negative numbers
	// such as in "-1 + 3" or in "1 + func(-10)" or again "1 + (-3)"
	outTree := Tree{}

	for i := 0; i < len(tree); i++ {
		e := tree[i]

		switch e.kind() {
		case treeEntryKind:
			subtree := e.(Tree).PrioritiseOperators()
			outTree = append(outTree, subtree)
			continue

		case operatorEntryKind:
			currentOperatorPrecedence := operatorPrecedence(e.(Operator))

			nextOperator := invalidOperator

			// TODO: check that i+1 is not out of range
			// TODO: does not (should it?) support expressions such as: 1+2+3 4*5
			for _, e2 := range tree[i+1:] {
				if e2.kind() == operatorEntryKind {
					nextOperator = e2.(Operator)
					break
				}
			}

			outTree = append(outTree, e)

			if nextOperator == invalidOperator || operatorPrecedence(nextOperator) <= currentOperatorPrecedence {
				continue
			}

			subTree := Tree{}

			// TODO: check that i+1 is not out of range
			// TODO: does not (should it?) support expressions such as: 1+2+3 4*5
			for _, e2 := range tree[i+1:] {
				if e2.kind() == operatorEntryKind && operatorPrecedence(e2.(Operator)) != operatorPrecedence(nextOperator) {
					break
				}
				subTree = append(subTree, e2)
				i++
			}

			outTree = append(outTree, subTree)

			i++
			outTree = append(outTree, tree[i])

		default:
			outTree = append(outTree, e)
			continue
		}
	}

	return outTree
}

func (Tree) kind() entryKind {
	return treeEntryKind
}
