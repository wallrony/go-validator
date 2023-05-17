package rules

import (
	"fmt"
	"reflect"
)

func NewRequiredRule(argument string) Rule {
	return &rule{
		typeName:    REQUIRED,
		description: "verify if a value exists",
		validator:   validateExists,
		argument:    argument,
	}
}

func validateExists(value interface{}) bool {
	if value == nil {
		return false
	}
	fieldType := reflect.TypeOf(value)
	switch fieldType.Kind() {
	case reflect.Slice:
		return reflect.ValueOf(value).Len() > 0
	case reflect.Map:
		return true
	case reflect.String:
		return value.(string) != ""
	}
	return true
}

func newRequiredError(fieldName, fieldType string) FieldError {
	message := fmt.Sprintf("'%s' field of type '%s' is missing or empty", fieldName, fieldType)
	return newFieldError(fieldName, message, REQUIRED)
}
