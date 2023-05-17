package rules

type validatorFunc func(value interface{}) bool

type Rule interface {
	Type() string
	Description() string
	Validator() validatorFunc
	Argument() string
	IsValid(value interface{}) bool
	GenerateError(fieldName string) FieldError
	IsSliceRule() bool
}

type rule struct {
	typeName    string
	description string
	validator   validatorFunc
	argument    string
}

const (
	REQUIRED         = "required"
	TYPE             = "type"
	LENGTH           = "length"
	MIN_LENGTH       = "minlength"
	MAX_LENGTH       = "maxlength"
	EMAIL_VALIDATION = "email"
	DATE_VALIDATION  = "date"
	ARRAY_LEN        = "slice:len"
	ARRAY_MIN_LEN    = "slice:minlen"
	ARRAY_MAX_LEN    = "slice:maxlen"
)

func (r *rule) Type() string {
	return r.typeName
}

func (r *rule) Description() string {
	return r.description
}

func (r *rule) Validator() validatorFunc {
	return r.validator
}

func (r *rule) Argument() string {
	return r.argument
}

func (r *rule) IsValid(value interface{}) bool {
	return r.validator(value)
}

func (r *rule) GenerateError(fieldName string) FieldError {
	return NewErrorByField(r.typeName, fieldName, r.argument)
}

func (r *rule) IsSliceRule() bool {
	return r.typeName == ARRAY_LEN || r.typeName == ARRAY_MIN_LEN || r.typeName == ARRAY_MAX_LEN
}
