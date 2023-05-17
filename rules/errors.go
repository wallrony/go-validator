package rules

type FieldError interface {
	Name() string
	Message() string
	RuleType() string
}

type fieldError struct {
	name     string
	message  string
	ruleType string
}

func newFieldError(name, message, ruleType string) FieldError {
	return &fieldError{name, message, ruleType}
}

func (f *fieldError) Name() string {
	return f.name
}

func (f *fieldError) Message() string {
	return f.message
}

func (f *fieldError) RuleType() string {
	return f.ruleType
}

func newErrorByType(t, fieldName, argument string) FieldError {
	switch t {
	case REQUIRED:
		return newRequiredError(fieldName, argument)
	case TYPE:
		return newWrongTypeError(fieldName, argument)
	case LENGTH:
		return newLengthError(fieldName, argument)
	case MIN_LENGTH:
		return newMinLengthError(fieldName, argument)
	case MAX_LENGTH:
		return newMaxLengthError(fieldName, argument)
	case EMAIL_VALIDATION:
		return newEmailValidationError(fieldName)
	case DATE_VALIDATION:
		return newDateValidationError(fieldName, argument)
	case ARRAY_LEN:
		return newArrayLenError(fieldName, argument)
	case ARRAY_MAX_LEN:
		return newArrayMaxlenError(fieldName, argument)
	case ARRAY_MIN_LEN:
		return newArrayMinlenError(fieldName, argument)
	}
	return nil
}

func NewErrorByField(ruleType, fieldName, argument string) FieldError {
	return newErrorByType(ruleType, fieldName, argument)
}
