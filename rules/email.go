package rules

import (
	"fmt"
	"net/mail"
)

func validateEmailRule() Rule {
	return &rule{
		typeName:    EMAIL_VALIDATION,
		description: "verify if a value is a valid email",
		validator:   validateEmailFN,
	}
}

func validateEmailFN(value interface{}) bool {
	if email, ok := value.(string); !ok {
		return false
	} else if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}

func newEmailValidationError(fieldName string) FieldError {
	message := fmt.Sprintf("the value provided for the '%s' field isn't a valid email", fieldName)
	return newFieldError(fieldName, message, EMAIL_VALIDATION)
}
