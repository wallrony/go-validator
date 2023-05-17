package rules

import (
	"fmt"
)

func newMinLengthRule(length int) Rule {
	return &rule{
		typeName:    MIN_LENGTH,
		description: fmt.Sprintf("verify if a value has a minimum length of %d", length),
		validator:   minLengthValidatorFN(length),
		argument:    fmt.Sprint(length),
	}
}

func newMaxLengthRule(length int) Rule {
	return &rule{
		typeName:    MAX_LENGTH,
		description: fmt.Sprintf("verify if a value has a maximum length of %d", length),
		validator:   maxLengthValidatorFN(length),
		argument:    fmt.Sprint(length),
	}
}

func newLengthRule(length int) Rule {
	return &rule{
		typeName:    LENGTH,
		description: fmt.Sprintf("verify if a value has length equals to %d", length),
		validator:   lengthValidatorFN(length),
		argument:    fmt.Sprint(length),
	}
}

func newLengthError(fieldName string, length string) FieldError {
	message := fmt.Sprintf("'%s' field must have %s characters", fieldName, length)
	return newFieldError(fieldName, message, LENGTH)
}

func newMinLengthError(fieldName string, length string) FieldError {
	message := fmt.Sprintf("'%s' field must have at least %s characters", fieldName, length)
	return newFieldError(fieldName, message, MIN_LENGTH)
}

func newMaxLengthError(fieldName string, length string) FieldError {
	message := fmt.Sprintf("'%s' field must have %s characters at max", fieldName, length)
	return newFieldError(fieldName, message, MAX_LENGTH)
}

func minLengthValidatorFN(length int) validatorFunc {
	return func(value interface{}) bool {
		if v, ok := value.(string); !ok {
			return false
		} else {
			return len(v) >= length
		}
	}
}

func maxLengthValidatorFN(length int) validatorFunc {
	return func(value interface{}) bool {
		if v, ok := value.(string); !ok {
			return false
		} else {
			return len(v) <= length
		}
	}
}

func lengthValidatorFN(length int) validatorFunc {
	return func(value interface{}) bool {
		if v, ok := value.(string); !ok {
			return false
		} else {
			return len(v) == length
		}
	}
}
