package rules

import (
	"fmt"
	"reflect"
)

// ArrayLen
func newArrayLenRule(len int) Rule {
	return &rule{
		typeName:    ARRAY_LEN,
		description: "verify if a value is an array and if has N elements",
		validator:   validateArrayLenFN(len),
		argument:    fmt.Sprint(len),
	}
}

func newArrayLenError(fieldName, argument string) FieldError {
	message := fmt.Sprintf("the '%s' field must have %s elements", fieldName, argument)
	return newFieldError(fieldName, message, ARRAY_LEN)
}

func validateArrayLenFN(len int) validatorFunc {
	return func(value interface{}) bool {
		if value == nil {
			return false
		} else if reflect.TypeOf(value).Kind() != reflect.Slice {
			return false
		} else {
			return reflect.ValueOf(value).Len() == len
		}
	}
}

// ArrayMaxlen
func newArrayMaxlenRule(maxlen int) Rule {
	return &rule{
		typeName:    ARRAY_MAX_LEN,
		description: "verify if a value is an array and if has N elements at max",
		validator:   validateArrayMaxlenFN(maxlen),
		argument:    fmt.Sprint(maxlen),
	}
}

func newArrayMaxlenError(fieldName, argument string) FieldError {
	message := fmt.Sprintf("the '%s' field must have %s elements at max", fieldName, argument)
	return newFieldError(fieldName, message, ARRAY_MAX_LEN)
}

func validateArrayMaxlenFN(maxlen int) validatorFunc {
	return func(value interface{}) bool {
		if value == nil {
			return false
		} else if reflect.TypeOf(value).Kind() != reflect.Slice {
			return false
		} else {
			return reflect.ValueOf(value).Len() <= maxlen
		}
	}
}

// ArrayMinlen
func newArrayMinlenRule(minlen int) Rule {
	return &rule{
		typeName:    ARRAY_MIN_LEN,
		description: "verify if a value is an array and if has at least N elements",
		validator:   validateArrayMinlenFN(minlen),
		argument:    fmt.Sprint(minlen),
	}
}

func newArrayMinlenError(fieldName, argument string) FieldError {
	message := fmt.Sprintf("the '%s' field must have at least %s elements", fieldName, argument)
	return newFieldError(fieldName, message, ARRAY_MIN_LEN)
}

func validateArrayMinlenFN(minlen int) validatorFunc {
	return func(value interface{}) bool {
		if value == nil {
			return false
		} else if reflect.TypeOf(value).Kind() != reflect.Slice {
			return false
		} else {
			return reflect.ValueOf(value).Len() >= minlen
		}
	}
}
