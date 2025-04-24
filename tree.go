package gal

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/samber/lo"
)

type entryKind int

func (ek entryKind) String() string {
	switch ek {
	case unknownEntryKind:
		return "unknownEntryKind"
	case valueEntryKind:
		return "valueEntryKind"
	case operatorEntryKind:
		return "operatorEntryKind"
	case treeEntryKind:
		return "treeEntryKind"
	case functionEntryKind:
		return "functionEntryKind"
	case objectPropertyEntryKind:
		return "objectPropertyEntryKind"
	case variableEntryKind:
		return "variableEntryKind"
	case objectAccessorEntryKind:
		return "objectAccessorEntryKind"
	default:
		return fmt.Sprintf("unknown:%d", ek)
	}
}

const (
	unknownEntryKind entryKind = iota
	valueEntryKind
	operatorEntryKind
	treeEntryKind
	functionEntryKind
	variableEntryKind
	objectPropertyEntryKind // "cousin" of a variableEntryKind, but for object properties
	objectAccessorEntryKind
)

type entry interface {
	kind() entryKind
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

// Variables holds the value of user-defined variables.
type Variables map[string]Value

func (v Variables) Get(name string) (Value, bool) {
	if v == nil {
		return nil, false
	}
	obj, ok := v[name]
	return obj, ok
}

// Functions holds the definition of user-defined functions.
type Functions map[string]FunctionalValue

func (f Functions) Get(name string) (FunctionalValue, bool) {
	if f == nil {
		return nil, false
	}
	obj, ok := f[name]
	return obj, ok
}

// Function returns the function definition of the function of the specified name.
// This method is used to look up object methods and user-defined functions.
// Built-in functions are not looked up here, they are pre-populated at
// parsing time by the TreeBuilder.
func (tc treeConfig) Function(name string) FunctionalValue {
	splits := strings.Split(name, ".")
	if len(splits) == 2 {
		// look up the method in the user provided objects
		if obj, ok := tc.objects.Get(splits[0]); ok {
			fv, _ := ObjectGetMethod(obj, splits[1])
			return fv
		}
		return func(...Value) Value {
			return NewUndefinedWithReasonf("error: object reference '%s' is not valid: unknown object or unknown method", name)
		}
	}

	if len(splits) >= 2 {
		return func(...Value) Value {
			return NewUndefinedWithReasonf("syntax error: object reference '%s' is not valid: too many dot accessors: max 1 permitted", name)
		}
	}

	// look up the function in the user-defined functions
	if val, ok := tc.functions.Get(name); ok {
		return val
	}

	return func(...Value) Value {
		return NewUndefinedWithReasonf("error: unknown user-defined function '%s'", name)
	}
}

// Objects is a collection of Object's in the form of a map which keys are the name of the
// object and values are the actual Object's.
type Objects map[string]Object

// Get returns the Object of the specified name.
func (o Objects) Get(name string) (Object, bool) {
	if o == nil {
		return nil, false
	}
	obj, ok := o[name]
	return obj, ok
}

type treeConfig struct {
	variables Variables
	functions Functions
	objects   Objects
}

// Variable returns the value of the variable specified by name.
// TODO: add support for arrays and maps via `[...]`
// ...   NOTE: it may be more adequate to create a new `[]` operator.
// ...   This would also permit its use on any Value, including those returned from function calls.
// ...   We would likely need to create new types (unless MultiValue can work for this).
// ...   An awkward and visually less elegant option would be builtin functions such as GetIndex() (for arrays) and GetKey (for maps).
// ...................................................................
// ...................................................................
// ...   Perhaps this indicates that it's time to drop gal.Value   ...
// ...   and use native Go types and reflection?!?!                ...
// ...................................................................
// ...................................................................
func (tc treeConfig) Variable(name string) Value {
	if val, ok := tc.variables.Get(name); ok {
		return val
	}

	return NewUndefinedWithReasonf("error: unknown user-defined variable '%s'", name)
}

func (tc treeConfig) ObjectProperty(objProp ObjectProperty) Value {
	if obj, ok := tc.objects.Get(objProp.ObjectName); ok {
		return ObjectGetProperty(obj, objProp.PropertyName)
	}
	return NewUndefinedWithReasonf("error: object property '%s': unknown object", objProp.String())
}

type treeOption func(*treeConfig)

// WithVariables is a functional parameter for Tree evaluation.
// It provides user-defined variables.
func WithVariables(vars Variables) treeOption {
	return func(cfg *treeConfig) {
		cfg.variables = vars
	}
}

// WithFunctions is a functional parameter for Tree evaluation.
// It provides user-defined functions.
func WithFunctions(funcs Functions) treeOption {
	return func(cfg *treeConfig) {
		cfg.functions = funcs
	}
}

// WithObjects is a functional parameter for Tree evaluation.
// It provides user-defined Objects.
// These objects can carry both properties and methods that can be accessed
// by gal in place of variables and functions.
func WithObjects(objects Objects) treeOption {
	return func(cfg *treeConfig) {
		cfg.objects = objects
	}
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
//
//nolint:maintidx
func (tree Tree) Calc(isOperatorInPrecedenceGroup func(Operator) bool, cfg *treeConfig) Tree {
	var (
		outTree Tree
		val     entry
		op      = invalidOperator
	)

	slog.Debug("Tree.Calc: start walking Tree", "tree", tree.String())

	//nolint:errcheck // life's too short to check for type assertion success here
	for i := 0; i < tree.TrunkLen(); i++ {
		if v, ok := val.(Undefined); ok {
			slog.Debug("Tree.Calc: val is Undefined", "i", i, "val", v.String())
			return Tree{v}
		}

		e := tree[i]
		slog.Debug("Tree.Calc: entry in Tree", "i", i, "kind", e.kind().String())
		if e == nil {
			slog.Debug("Tree.Calc: nil entry in Tree")
			return Tree{
				NewUndefinedWithReasonf("syntax error: nil value at tree entry #%d - tree: %+v", i, tree),
			}
		}

		switch e.kind() {
		case valueEntryKind:
			slog.Debug("Tree.Calc: valueEntryKind", "i", i, "Value", e.(Value).String())
			if val == nil && op == invalidOperator {
				val = e
				continue
			}

			if val == nil {
				return Tree{
					NewUndefinedWithReasonf("syntax error: missing left hand side value for operator '%s'", op.String()),
				}
			}

			slog.Debug("Tree.Calc: valueEntryKind - calculate", "i", i, "val", val.(Value).String(), "op", op.String(), "e", e.(Value).String())
			val = calculate(val.(Value), op, e.(Value))
			slog.Debug("Tree.Calc: valueEntryKind - calculate", "i", i, "result", val.(Value).String())

		case treeEntryKind:
			slog.Debug("Tree.Calc: treeEntryKind", "i", i)
			if val == nil && op != invalidOperator {
				return Tree{
					NewUndefinedWithReasonf("syntax error: missing left hand side value for operator '%s'", op.String()),
				}
			}

			rhsVal := e.(Tree).Eval(WithFunctions(cfg.functions), WithVariables(cfg.variables), WithObjects(cfg.objects))
			if v, ok := rhsVal.(Undefined); ok {
				slog.Debug("Tree.Calc: val is Undefined", "i", i, "val", v.String())
				return Tree{v}
			}
			if val == nil {
				val = rhsVal
				continue
			}

			val = calculate(val.(Value), op, rhsVal)
			slog.Debug("Tree.Calc: treeEntryKind - calculate", "i", i, "val", val.(Value).String(), "op", op.String(), "rhsVal", rhsVal.String(), "result", val.(Value).String())

		case operatorEntryKind:
			slog.Debug("Tree.Calc: operatorEntryKind", "i", i, "Value", e.(Operator).String())
			op = e.(Operator) //nolint:errcheck
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

		case functionEntryKind:
			slog.Debug("Tree.Calc: functionEntryKind", "i", i, "name", e.(Function).Name)
			f := e.(Function) //nolint:errcheck
			if f.BodyFn == nil {
				// attempt to get body of a user-defined variable or a user-provided object's method.
				f.BodyFn = cfg.Function(f.Name)
			}

			rhsVal := f.Eval(WithFunctions(cfg.functions), WithVariables(cfg.variables), WithObjects(cfg.objects))
			if v, ok := rhsVal.(Undefined); ok {
				slog.Debug("Tree.Calc: val is Undefined", "i", i, "val", v.String())
				return Tree{v}
			}
			if val == nil {
				val = rhsVal
				continue
			}

			lhsVal := val
			val = calculate(val.(Value), op, rhsVal)
			slog.Debug("Tree.Calc: functionEntryKind - calculate", "i", i, "lhsVal", lhsVal.(Value).String(), "op", op.String(), "rhsVal", rhsVal.String(), "result", val.(Value).String())

		case variableEntryKind:
			val = variableEntryKindFn(i, val, op, e, cfg)

		case objectPropertyEntryKind:
			slog.Debug("Tree.Calc: objectPropertyEntryKind", "i", i, "object_property", e.(ObjectProperty).String())
			objProp := e.(ObjectProperty)
			rhsVal := cfg.ObjectProperty(objProp)
			slog.Debug("Tree.Calc: objectPropertyEntryKind", "i", i, "value", rhsVal.String())

			if val == nil {
				val = rhsVal
				continue
			}

			val = calculate(val.(Value), op, rhsVal)
			slog.Debug("Tree.Calc: objectPropertyEntryKind - calculate", "i", i, "val", val.(Value).String(), "op", op.String(), "rhsVal", rhsVal.String(), "result", val.(Value).String())

		case objectAccessorEntryKind:
			switch a := e.(type) {
			case Dot[Function]:
				slog.Debug("Tree.Calc: objectAccessorEntryKind Dot[Function]", "i", i, "member_name", a.Member.Name)
				f := a.Member
				if f.BodyFn != nil {
					return Tree{
						// NOTE: this could be supported but it would turn the object into a prototype model e.g. like JavaScript
						NewUndefinedWithReasonf("internal error: objectAccessorEntryKind Dot[Function] for '%s': BodyFn is not empty: this indicates the object's method was confused for a build-in function", f.Name),
					}
				}
				// as this is an object function accessor, we need to get the object first: it is the LHS currently held in val
				vVal, ok := val.(Value)
				if !ok {
					return Tree{
						NewUndefinedWithReasonf("syntax error: object accessor called on non-object: [object: '%s'] [member: '%s']", val.kind().String(), f.Name),
					}
				}
				// now, we can get the method from the object
				if vFv, ok := ObjectGetMethod(vVal, f.Name); ok {
					f.BodyFn = vFv
					rhsVal := f.Eval(WithFunctions(cfg.functions), WithVariables(cfg.variables), WithObjects(cfg.objects))
					if v, ok := rhsVal.(Undefined); ok {
						slog.Debug("Tree.Calc: object method's BodyFn is Undefined", "i", i, "val", v.String())
						return Tree{v}
					}
					val = rhsVal
					continue
				} else {
					return Tree{
						NewUndefinedWithReasonf("syntax error: object accessor function called on unknown or non-function member: [object: '%s'] [member: '%s']", val.kind().String(), f.Name),
					}
				}

			case Dot[Variable]:
				slog.Debug("Tree.Calc: objectAccessorEntryKind Dot[Variable]", "i", i, "member_name", a.Member.Name)
				v := a.Member
				// as this is an object property accessor, we need to get the object first: it is the LHS currently held in val
				var vVal any
				vVal, ok := val.(Value)
				if !ok {
					return Tree{
						NewUndefinedWithReasonf("syntax error: object accessor called on non-object: [object: '%s'] [member: '%s']", val.kind().String(), v.Name),
					}
				}
				// if the object is a ObjectValue, we need to get the underlying object
				// ObjectValue is a wrapper for "general" objects (i.e. non-gal.Value objects)
				// By Object, we mean a Go struct, a pointer to a struct or a Go interface.
				objVal, ok := vVal.(ObjectValue)
				if ok {
					vVal = objVal.Object
				}
				// now, we can get the property from the object
				vFv := ObjectGetProperty(vVal, v.Name)
				rhsVal := vFv
				if v, ok := rhsVal.(Undefined); ok {
					slog.Debug("Tree.Calc: object property's value is Undefined", "i", i, "val", v.String())
					return Tree{v}
				}
				val = rhsVal
				continue

			default:
				slog.Debug("Tree.Calc: objectAccessorEntryKind Dot[unknown]", "i", i, "entry_string", a.kind().String())
				return Tree{
					NewUndefinedWithReasonf("internal error: unknown objectAccessorEntryKind Dot kind: '%s'", e.kind().String()),
				}
			}

		case unknownEntryKind:
			slog.Debug("Tree.Calc: unknownEntryKind", "i", i, "val", val, "op", op.String(), "e", e)
			return Tree{e}

		default:
			slog.Debug("Tree.Calc: default case", "i", i, "val", val, "op", op.String(), "e", e)
			return Tree{
				NewUndefinedWithReasonf("internal error: unknown entry kind: '%s'", e.kind().String()),
			}
		}
	}

	if val != nil {
		outTree = append(outTree, val)
	}

	return outTree
}

//nolint:errcheck // life's too short to check for type assertion success here
func variableEntryKindFn(i int, val entry, op Operator, e entry, cfg *treeConfig) entry {
	varName := e.(Variable).Name
	slog.Debug("Tree.Calc: variableEntryKind", "i", i, "name", varName)

	rhsVal := cfg.Variable(varName)
	slog.Debug("Tree.Calc: variableEntryKind", "i", i, "value", rhsVal.String())

	if val == nil {
		val = rhsVal
		return val
	}

	val = calculate(val.(Value), op, rhsVal)
	slog.Debug("Tree.Calc: variableEntryKind - calculate", "i", i, "val", val.(Value).String(), "op", op.String(), "rhsVal", rhsVal.String(), "result", val.(Value).String())

	return val
}

// CleanUp performs simplification operations before calculating this tree.
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
		return append(Tree{NewNumberFromInt(-1), Multiply}, outTree[1:]...)
	}

	panic("point never reached")
}

func (Tree) kind() entryKind {
	return treeEntryKind
}

func (tree Tree) String(indents ...string) string {
	indent := strings.Join(indents, "")

	res := ""
	for _, e := range tree {
		//nolint:errcheck // life's too short to check for type assertion success here
		switch e.kind() {
		case unknownEntryKind:
			res += fmt.Sprintf(indent+"unknownEntryKind %T\n", e)
		case valueEntryKind:
			res += fmt.Sprintf(indent+"Value %T %s\n", e, e.(Value).String())
		case operatorEntryKind:
			res += fmt.Sprintf(indent+"Operator %s\n", e.(Operator).String())
		case treeEntryKind:
			res += fmt.Sprintf(indent+"Tree {\n%s}\n", e.(Tree).String("   "))
		case functionEntryKind:
			args := lo.Map(e.(Function).Args, func(item Tree, index int) string {
				return strings.TrimRight(item.String(), "\n")
			})
			res += fmt.Sprintf(indent+"Function %s(%s)\n", e.(Function).Name, strings.Join(args, ", "))
		case variableEntryKind:
			res += fmt.Sprintf(indent+"Variable %s\n", e.(Variable).Name)
		case objectPropertyEntryKind:
			res += fmt.Sprintf(indent+"ObjectProperty %s.%s\n", e.(ObjectProperty).ObjectName, e.(ObjectProperty).PropertyName)
		case objectAccessorEntryKind:
			switch a := e.(type) {
			case Dot[Function]:
				fn := a.Member
				args := lo.Map(fn.Args, func(item Tree, index int) string {
					return strings.TrimRight(item.String(), "\n")
				})
				res += fmt.Sprintf(indent+"ObjectAccessor %s(%s)\n", fn.Name, strings.Join(args, ", "))
			case Dot[Variable]:
				v := a.Member
				res += fmt.Sprintf(indent+"ObjectAccessor %s\n", v.Name)
			default:
				res += fmt.Sprintf(indent+"TODO: unsupported - %s %T %s\n", e, a, a.kind().String())
			}
		default:
			res += fmt.Sprintf(indent+"TODO: unsupported - %T %s\n", e, e.kind().String())
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
