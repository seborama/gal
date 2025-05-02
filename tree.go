package gal

import (
	"fmt"
	"strings"
)

type entry interface {
	// TODO: remove this interface and perhaps use Value instead? Or else, use 'any'
}

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
		if v, ok := val.(Undefined); ok {
			// TODO: is this odd? we also perform this check below, after the switch case
			return Tree{v}
		}

		e := tree[i]
		if e == nil {
			return Tree{
				NewUndefinedWithReasonf("syntax error: nil value at tree entry #%d - tree: %+v", i, tree),
			}
		}

		switch typedE := e.(type) {
		case Bool, MultiValue, Number, String:
			// TODO: (!!) implement Calc() on all these types (should be possible to encapsulate the behaviour to avoid repeat code)
			fmt.Printf("DEBUG - Tree.Calc: entry in Tree %T\n", typedE)
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

		case Dot[Function]: // TODO: (!!) it doesn't seem to make sense to have Dot[Member] anymore since the implementations of Calculate diverge.
			val = objectAccessorDotFunctionFn(val, typedE, cfg)

		case Dot[Variable]:
			val = objectAccessorDotVariableFn(val, typedE)

		case Undefined:
			return Tree{e}

		default:
			val = NewUndefinedWithReasonf("internal error: unknown entry type: '%T'", e)
		}

		_, isUndef := val.(Undefined)
		if val != nil && isUndef {
			// TODO: (!!) is this odd? we also perform this check above, at the start of te for loop
			return Tree{val}
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
	if v, ok := rhsVal.(Undefined); ok {
		return v
	}

	if val == nil {
		return rhsVal
	}

	val = calculate(val.(Value), op, rhsVal)

	return val
}

//nolint:errcheck // life's too short to check for type assertion success here
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

func objectAccessorDotFunctionFn(val entry, oa Dot[Function], cfg *treeConfig) entry {
	fn := oa.Member

	if fn.BodyFn != nil {
		// NOTE: this could be supported but it would turn the object into a prototype model e.g. like JavaScript
		return NewUndefinedWithReasonf("internal error: objectAccessorEntryKind Dot[Function] for '%s': BodyFn is not empty: this indicates the object's method was confused for a build-in function", fn.Name)
	}

	// as this is an object function accessor, we need to get the object first: it is the LHS currently held in val
	vVal, ok := val.(Value)
	if !ok {
		return NewUndefinedWithReasonf("syntax error: object accessor called on non-object: [object: '%T'] [member: '%s']", val, fn.Name)
	}

	// now, we can get the method from the object
	if vFv, ok := ObjectGetMethod(vVal, fn.Name); ok {
		fn.BodyFn = vFv
		rhsVal := fn.Eval(WithFunctions(cfg.functions), WithVariables(cfg.variables), WithObjects(cfg.objects))
		if v, ok := rhsVal.(Undefined); ok {
			return v
		}

		return rhsVal
	}

	return NewUndefinedWithReasonf("syntax error: object accessor function called on unknown or non-function member: [object: '%T'] [member: '%s']", fmt.Sprintf("%T", val), fn.Name)
}

func objectAccessorDotVariableFn(val entry, oa Dot[Variable]) entry {
	v := oa.Member

	// as this is an object property accessor, we need to get the object first: it is the LHS currently held in val
	var vVal any
	vVal, ok := val.(Value)
	if !ok {
		return NewUndefinedWithReasonf("syntax error: object accessor called on non-object: [object: '%T'] [member: '%s']", fmt.Sprintf("%T", val), v.Name)
	}

	// if the object is a ObjectValue, we need to get the underlying object
	// ObjectValue is a wrapper for "general" objects (i.e. non-gal.Value objects)
	// By Object, we mean a Go struct, a pointer to a struct or a Go interface.
	objVal, ok := vVal.(ObjectValue)
	if ok {
		vVal = objVal.Object
	}

	// now, we can get the property from the object
	rhsVal := ObjectGetProperty(vVal, v.Name)

	return rhsVal
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
		//nolint:errcheck // life's too short to check for type assertion success here
		switch typedE := e.(type) {
		case Undefined:
			res += fmt.Sprintf(indent+"unknownEntryKind %T\n", e)
		case Value:
			res += fmt.Sprintf(indent+"Value %T %s\n", e, typedE.String())
		case Operator:
			res += fmt.Sprintf(indent+"Operator %s\n", typedE.String())
		case Tree:
			res += fmt.Sprintf(indent+"Tree {\n%s}\n", typedE.String("   "))
		case Function:
			res += fmt.Sprintf(indent+"Function %s(%s)\n", typedE.String())
		case Variable:
			res += fmt.Sprintf(indent+"Variable %s\n", typedE.Name)
		case ObjectProperty:
			res += fmt.Sprintf(indent+"ObjectProperty %s\n", typedE.String())
		case ObjectMethod:
			res += fmt.Sprintf(indent+"ObjectMethod %s\n", typedE.String())
		case Dot[Function], Dot[Variable]: // TODO: split this into two cases
			switch a := e.(type) {
			case Dot[Function]:
				fn := a.Member
				res += fmt.Sprintf(indent+"ObjectAccessor[Function] %s\n", fn.String())
			case Dot[Variable]:
				v := a.Member
				res += fmt.Sprintf(indent+"ObjectAccessor[Variable] %s\n", v.String())
			default:
				res += fmt.Sprintf(indent+"TODO: unsupported - %s %T\n", e, a) // TODO: does %s on e work here????
			}
		default:
			res += fmt.Sprintf(indent+"TODO: unsupported - %T\n", e)
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
