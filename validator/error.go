package validator

import (
	"go-validator/rules"
	"strings"
)

type ValidationError interface {
	String() string
	Messages() []string
	Fields() []string
	RuleTypes() []string
	FieldsErrors() []rules.FieldError
}

type validationError struct {
	fieldErrors []rules.FieldError
}

func newValidationError(fieldErrors []rules.FieldError) ValidationError {
	return &validationError{fieldErrors}
}

func (v *validationError) String() string {
	return strings.Join(v.Messages(), " & ")
}

func (v *validationError) Messages() []string {
	var messages []string
	for _, fieldError := range v.fieldErrors {
		messages = append(messages, fieldError.Message())
	}
	return messages
}

func (v *validationError) Fields() []string {
	var fields []string
	for _, fieldError := range v.fieldErrors {
		fields = append(fields, fieldError.Name())
	}
	return fields
}

func (v *validationError) RuleTypes() []string {
	var ruleTypes []string
	for _, fieldError := range v.fieldErrors {
		ruleTypes = append(ruleTypes, fieldError.RuleType())
	}
	return ruleTypes
}

func (v *validationError) FieldsErrors() []rules.FieldError {
	return v.fieldErrors
}
