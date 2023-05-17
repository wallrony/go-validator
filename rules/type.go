package rules

import (
	"fmt"
	"reflect"
)

var typeValidator = map[string]Rule{
	"int":     NewTypeRule[int](),
	"int32":   NewTypeRule[float64](),
	"int64":   NewTypeRule[float64](),
	"float32": NewTypeRule[float64](),
	"float64": NewTypeRule[float64](),
	"string":  NewTypeRule[string](),
	"bool":    NewTypeRule[bool](),
}

func NewTypeRuleWithMethod(method validatorFunc) Rule {
	return &rule{
		typeName:    TYPE,
		description: "verify if a value is convertable to int",
		validator:   method,
		argument:    "int",
	}
}

func GetType[T comparable]() string {
	var t T
	return fmt.Sprintf("%T", t)
}

func NewTypeRule[T comparable]() Rule {
	typeName := GetType[T]()
	return &rule{
		typeName:    TYPE,
		description: fmt.Sprintf("verify if a value is convertable to %s", typeName),
		validator:   validateType[T],
		argument:    typeName,
	}
}

func NewSliceRule[T comparable]() Rule {
	typeName := GetType[T]()
	return &rule{
		typeName:    TYPE,
		description: fmt.Sprintf("verify if a value is convertable to %s array", typeName),
		validator:   validateSliceType[T],
		argument:    typeName,
	}
}

func GetTypeValidator(key string) Rule {
	return typeValidator[key]
}

func newWrongTypeError(fieldName, fieldType string) FieldError {
	if fieldType == "struct" {
		fieldType = "json"
	}
	message := fmt.Sprintf("'%s' field type must be '%s'", fieldName, fieldType)
	return newFieldError(fieldName, message, TYPE)
}

func validateType[T comparable](value interface{}) bool {
	if value == nil {
		return true
	}
	var t T
	neededType := reflect.TypeOf(t)
	valueType := reflect.TypeOf(value)
	return neededType == valueType || valueType.ConvertibleTo(neededType)
}

func validateSliceType[T comparable](value interface{}) bool {
	newValue, ok := value.([]interface{})
	for _, item := range newValue {
		if _, ok = item.(T); !ok {
			return false
		}
	}
	return ok
}
