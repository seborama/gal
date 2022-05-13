package gal

type entryKind int

const (
	unknownEntryKind entryKind = iota
	valueEntryKind
	operatorEntryKind
	treeEntryKind
	functionEntryKind
)

type entry interface {
	kind() entryKind
}

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

// Split divides a Tree trunk at points where twoconsecutive Values are present.
func (tree Tree) Split() []Tree {
	var forest []Tree

	partStart := 0

	for i := 1; i < tree.TrunkLen(); i++ {
		_, ok1 := tree[i].(Value)
		_, ok2 := tree[i-1].(Value)

		if ok1 && ok2 {
			forest = append(forest, tree[partStart:i])
			partStart = i
		}
	}

	return append(forest, tree[partStart:])
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
					NewUndefinedWithReasonf("syntax error: nil value cannot be operated upon (op='%s')", op.String()),
				}
			}

			val = calculate(val.(Value), op, e.(Value))

		case treeEntryKind:
			if val == nil && op != invalidOperator {
				return Tree{
					NewUndefinedWithReasonf("syntax error: nil value cannot be operated upon (op='%s')", op.String()),
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

		case functionEntryKind:
			rhsVal := e.(Function).Eval()
			if val == nil {
				val = rhsVal
				continue
			}

			val = calculate(val.(Value), op, rhsVal)

		case unknownEntryKind:
			// TODO: distinguish between unknownEntryKind and undefinedEntryKind (which is a Value)
			return Tree{e}

		default:
			return Tree{
				NewUndefinedWithReasonf("internal error: unknown entry kind: '%v'", e.kind()),
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
	case Plus:
		outVal = lhs.Add(rhs)

	case Minus:
		outVal = lhs.Sub(rhs)

	case Multiply:
		outVal = lhs.Multiply(rhs)

	case Divide:
		outVal = lhs.Divide(rhs)

	case Power:
		outVal = lhs.PowerOf(rhs)

	case Modulus:
		outVal = lhs.Mod(rhs)

	default:
		return NewUndefinedWithReasonf("unimplemented operator: '%s'", op.String())
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
			case Plus:
				// drop superfluous plus sign at start of Tree
				continue

			case Minus:
				outTree = append(outTree, NewNumber(-1), Multiply, tree[i+1])
				i++
				continue

			default:
				return Tree{
					NewUndefinedWithReasonf("syntax error: expression starts with '%s'", e.(Operator).String()),
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
