package rules

import (
	"fmt"
	"time"
)

func validateDateRule(format string) Rule {
	return &rule{
		typeName:    DATE_VALIDATION,
		description: "verify if a value is a valid date",
		validator:   validateDateFN(format),
		argument:    format,
	}
}

func validateDateFN(format string) validatorFunc {
	return func(value interface{}) bool {
		if date, ok := value.(string); !ok {
			return false
		} else if _, err := time.Parse(format, date); err != nil {
			return false
		}
		return true
	}
}

func newDateValidationError(fieldName, format string) FieldError {
	message := fmt.Sprintf("'%s' field doesn't match with the '%s' format", fieldName, format)
	return newFieldError(fieldName, message, DATE_VALIDATION)
}
