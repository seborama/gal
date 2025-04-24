package gal

import "fmt"

type Variable struct {
	Name string
}

func NewVariable(name string) Variable {
	return Variable{
		Name: name,
	}
}

func (Variable) kind() entryKind {
	return variableEntryKind
}

func (v Variable) String() string {
	return v.Name
}

type ObjectProperty struct {
	ObjectName   string
	PropertyName string
}

func NewObjectProperty(objectName, propertyName string) ObjectProperty {
	return ObjectProperty{
		ObjectName:   objectName,
		PropertyName: propertyName,
	}
}

func (o ObjectProperty) kind() entryKind {
	return objectPropertyEntryKind
}

func (o ObjectProperty) String() string {
	return fmt.Sprintf("%s.%s", o.ObjectName, o.PropertyName)
}
