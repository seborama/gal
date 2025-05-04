package gal

import (
	"fmt"
	"strings"
)

type entry interface{} // NOTE: this could be dropped in favour of using `any` directly.

type Tree []entry

func (tree Tree) TrunkLen() int {
	return len(tree)
}

// FullLen returns the total number of non 'Tree-type' elements in the tree.
func (tree Tree) FullLen() int {
	l := len(tree)

	for _, e := range tree {
		if subTree, ok := e.(Tree); ok {
			l += subTree.FullLen() - 1
		}
	}

	return l
}

// Eval evaluates this tree and returns its value.
// It accepts optional functional parameters to supply user-defined
// entities such as functions and variables.
func (tree Tree) Eval(opts ...treeOption) Value {
	// config
	cfg := &treeConfig{}

	for _, o := range opts {
		o(cfg)
	}

	// Execute calculation by decreasing order of precedence.
	// It is necessary to proceed by operator precedence in order
	// to calculate the expression under conventional rules of precedence.
	workingTree := tree.
		CleanUp().
		Calc(powerOperators, cfg).
		Calc(multiplicativeOperators, cfg).
		Calc(additiveOperators, cfg).
		Calc(bitwiseShiftOperators, cfg).
		Calc(comparativeOperators, cfg).
		Calc(logicalOperators, cfg)

	// TODO: refactor this
	// perhaps add Tree.Value() which tests that only one entry is left and that it is a Value
	// (maybe MultiValue can help too?)
	//nolint:errcheck // life's too short to check for type assertion success here
	return workingTree[0].(Value)
}

// Split divides a Tree trunk at points where two consecutive entries are present without
// an operator in between.
func (tree Tree) Split() []Tree {
	if len(tree) == 0 {
		return []Tree{}
	}

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

// Calc is a reduction operation that calculates the Value of sub-expressions contained
// in this Tree, based on operator precedence.
// When isOperatorInPrecedenceGroup returns true, the operator is calculated and the resultant
// Value is inserted in _replacement_ of the terms (elements) of this Tree that where calculated.
// For instance, a tree representing the expression '2 + 5 * 4 / 2' with an operator precedence
// of 'multiplicativeOperators' would read the Tree left to right and return a new Tree that
// represents: '2 + 10' where 10 was calculated (and reduced) from 5 * 4 = 20 / 2 = 10.
func (tree Tree) Calc(isOperatorInPrecedenceGroup func(Operator) bool, cfg *treeConfig) Tree {
	var (
		outTree Tree
		val     entry
		op      = invalidOperator
	)

	//nolint:errcheck // life's too short to check for type assertion success here
	for i := 0; i < tree.TrunkLen(); i++ {
		if u, ok := val.(Undefined); ok {
			return Tree{u}
		}

		e := tree[i]
		if e == nil {
			return Tree{
				NewUndefinedWithReasonf("syntax error: nil value at tree entry #%d - tree: %+v", i, tree),
			}
		}

		switch typedE := e.(type) {
		case Bool, MultiValue, Number, String:
			vVal, _ := val.(Value) // avoid panic if val is nil
			val = valueEntryKindFn(vVal, op, e.(Value))

		case Tree:
			val = typedE.Calculate(val, op, cfg)

		case Operator:
			op = typedE
			if isOperatorInPrecedenceGroup(op) {
				// same operator precedence: keep operating linearly, do not build a tree
				continue
			}
			if val != nil {
				outTree = append(outTree, val)
			}
			outTree = append(outTree, op)
			// just found and processed the current operator - now, reset val and op and start again from fresh
			val = nil
			op = invalidOperator

		case Function:
			val = typedE.Calculate(val, op, cfg)

		case ObjectMethod:
			val = typedE.Calculate(val, op, cfg)

		case Variable:
			val = typedE.Calculate(val, op, cfg)

		case ObjectProperty:
			val = typedE.Calculate(val, op, cfg)

		case DotFunction:
			val = typedE.Calculate(val, cfg)

		case DotVariable:
			val = typedE.Calculate(val)

		case Undefined:
			return Tree{e}

		default:
			val = NewUndefinedWithReasonf("internal error: unknown entry type: '%T'", e)
		}
	}

	if val != nil {
		outTree = append(outTree, val)
	}

	return outTree
}

func (tree Tree) Calculate(val entry, op Operator, cfg *treeConfig) entry {
	if val == nil && op != invalidOperator {
		return NewUndefinedWithReasonf("syntax error: missing left hand side value for operator '%s'", op.String())
	}

	rhsVal := tree.Eval(WithFunctions(cfg.functions), WithVariables(cfg.variables), WithObjects(cfg.objects))
	if u, ok := rhsVal.(Undefined); ok {
		return u
	}

	if val == nil {
		return rhsVal
	}

	//nolint:errcheck // life's too short to check for type assertion success here
	val = calculate(val.(Value), op, rhsVal)

	return val
}

func valueEntryKindFn(val Value, op Operator, e Value) entry {
	if val == nil && op == invalidOperator {
		return e
	}

	if val == nil {
		return NewUndefinedWithReasonf("syntax error: missing left hand side value for operator '%s'", op.String())
	}

	val = calculate(val, op, e)

	return val
}

// CleanUp performs simplification operations before calculating this tree.
func (tree Tree) CleanUp() Tree {
	return tree.cleansePlusMinusTreeStart()
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
		return append(Tree{NewNumberFromInt(-1), Multiply}, outTree[1:]...)
	}

	panic("this point should never be reached")
}

func (tree Tree) String(indents ...string) string {
	indent := strings.Join(indents, "")

	res := ""
	for _, e := range tree {
		switch typedE := e.(type) {
		case Undefined:
			res += fmt.Sprintf("%sunknownEntryKind %T\n", indent, e)
		case Value:
			res += fmt.Sprintf("%sValue %T %s\n", indent, e, typedE.String())
		case Operator:
			res += fmt.Sprintf("%sOperator %s\n", indent, typedE.String())
		case Tree:
			res += fmt.Sprintf("%sTree {\n%s}\n", indent, typedE.String("   "))
		case Function:
			res += fmt.Sprintf("%sFunction %s\n", indent, typedE.String())
		case Variable:
			res += fmt.Sprintf("%sVariable %s\n", indent, typedE.Name)
		case ObjectProperty:
			res += fmt.Sprintf("%sObjectProperty %s\n", indent, typedE.String())
		case ObjectMethod:
			res += fmt.Sprintf("%sObjectMethod %s\n", indent, typedE.String())
		case DotFunction:
			res += fmt.Sprintf("%sDotFunction %s\n", indent, typedE.String())
		case DotVariable:
			res += fmt.Sprintf("%sDotVariable %s\n", indent, typedE.String())
		default:
			res += fmt.Sprintf("%sTODO: unsupported - %T\n", indent, e)
		}
	}

	return strings.TrimRight(res, "\n")
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

	case LessThan:
		outVal = lhs.LessThan(rhs)

	case LessThanOrEqual:
		outVal = lhs.LessThanOrEqual(rhs)

	case EqualTo:
		outVal = lhs.EqualTo(rhs)

	case NotEqualTo:
		outVal = lhs.NotEqualTo(rhs)

	case GreaterThan:
		outVal = lhs.GreaterThan(rhs)

	case GreaterThanOrEqual:
		outVal = lhs.GreaterThanOrEqual(rhs)

	case And, And2:
		outVal = lhs.And(rhs)

	case Or, Or2:
		outVal = lhs.Or(rhs)

	default:
		return NewUndefinedWithReasonf("unimplemented operator: '%s'", op.String())
	}

	return outVal
}
