package gal

type entryKind int

const (
	unknownEntryKind entryKind = iota
	valueEntryKind
	operatorEntryKind
	treeEntryKind
	functionEntryKind
	variableEntryKind
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

type Variables map[string]Value

type Functions map[string]FunctionalValue

func (f Functions) Function(name string) FunctionalValue {
	if f == nil {
		return nil
	}

	if val, ok := f[name]; ok {
		return val
	}

	return nil
}

type treeConfig struct {
	variables Variables
	functions Functions
}

func (tc treeConfig) Variable(name string) (Value, bool) {
	if tc.variables == nil {
		return nil, false
	}

	val, ok := tc.variables[name]
	return val, ok
}

type treeOption func(*treeConfig)

func WithVariables(vars Variables) treeOption {
	return func(cfg *treeConfig) {
		cfg.variables = vars
	}
}

func WithFunctions(funcs Functions) treeOption {
	return func(cfg *treeConfig) {
		cfg.functions = funcs
	}
}

func (tree Tree) Eval(opts ...treeOption) Value {
	//config
	cfg := &treeConfig{}

	for _, o := range opts {
		o(cfg)
	}

	// execute calculation by decreasing order of precedence
	workingTree := tree.
		CleanUp().
		Calc(powerOperators, cfg).
		Calc(multiplicativeOperators, cfg).
		Calc(additiveOperators, cfg).
		Calc(bitwiseShiftOperators, cfg)

	// TODO: refactor this
	// perhaps add Tree.Value() which tests that only one entry is left and that it is a Value
	return workingTree[0].(Value)
}

// Split divides a Tree trunk at points where two consecutive entries are present without
// an operator in between.
func (tree Tree) Split() []Tree {
	var forest []Tree

	partStart := 0

	for i := 1; i < tree.TrunkLen(); i++ {
		_, ok1 := tree[i].(Operator)
		_, ok2 := tree[i-1].(Operator)

		if !ok1 && !ok2 {
			forest = append(forest, tree[partStart:i])
			partStart = i
		}
	}

	return append(forest, tree[partStart:])
}

func (tree Tree) Calc(isOperatorInPrecedenceGroup func(Operator) bool, cfg *treeConfig) Tree {
	var outTree Tree

	var val entry
	var op Operator = invalidOperator

	for i := 0; i < tree.TrunkLen(); i++ {
		e := tree[i]
		if e == nil {
			return Tree{
				NewUndefinedWithReasonf("syntax error: nil value at tree entry #%d - tree: %+v", i, tree),
			}
		}

		switch e.kind() {
		case valueEntryKind:
			if val == nil && op == invalidOperator {
				val = e
				continue
			}

			if val == nil {
				return Tree{
					NewUndefinedWithReasonf("syntax error: missing left hand side value for operator '%s'", op.String()),
				}
			}

			val = calculate(val.(Value), op, e.(Value))

		case treeEntryKind:
			if val == nil && op != invalidOperator {
				return Tree{
					NewUndefinedWithReasonf("syntax error: missing left hand side value for operator '%s'", op.String()),
				}
			}

			rhsVal := e.(Tree).Eval(WithFunctions(cfg.functions), WithVariables(cfg.variables))
			if val == nil {
				val = rhsVal
				continue
			}

			val = calculate(val.(Value), op, rhsVal)

		case operatorEntryKind:
			op = e.(Operator)
			if isOperatorInPrecedenceGroup(op) {
				continue
			}
			if val != nil {
				outTree = append(outTree, val)
			}
			outTree = append(outTree, op)
			val = nil
			op = invalidOperator

		case functionEntryKind:
			f := e.(Function)
			if f.BodyFn == nil {
				f.BodyFn = cfg.functions.Function(f.Name)
			}
			rhsVal := f.Eval(WithFunctions(cfg.functions), WithVariables(cfg.variables))
			if val == nil {
				val = rhsVal
				continue
			}

			val = calculate(val.(Value), op, rhsVal)

		case variableEntryKind:
			varName := e.(Variable).Name
			rhsVal, ok := cfg.Variable(varName)
			if !ok {
				return Tree{
					NewUndefinedWithReasonf("syntax error: unknown variable name: '%s'", varName),
				}
			}

			if val == nil {
				val = rhsVal
				continue
			}

			val = calculate(val.(Value), op, rhsVal)

		case unknownEntryKind:
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
	var outVal Value

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

	case LShift:
		outVal = lhs.LShift(rhs)

	case RShift:
		outVal = lhs.RShift(rhs)

	default:
		return NewUndefinedWithReasonf("unimplemented operator: '%s'", op.String())
	}

	return outVal
}

func (tree Tree) CleanUp() Tree {
	return tree.
		cleansePlusMinusTreeStart()
}

// cleansePlusMinusTreeStart consolidates the - and + that are at the first position in a Tree.
// `plus` is removed and `minus` causes the number that follows to be negated.
func (tree Tree) cleansePlusMinusTreeStart() Tree {
	outTree := make(Tree, len(tree))
	copy(outTree, tree)

	if tree.TrunkLen() < 2 || (tree[0] != Plus && tree[0] != Minus) {
		return outTree
	}

	switch outTree[0] {
	case Plus:
		return outTree[1:]
	case Minus:
		return append(Tree{NewNumber(-1), Multiply}, outTree[1:]...)
	}

	panic("point never reached")
}

func (Tree) kind() entryKind {
	return treeEntryKind
}
