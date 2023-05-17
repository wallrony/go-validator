package validator

import (
	"fmt"
	"go-validator/rules"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	// Tags
	validationTag     = "validate"
	jsonTag           = "json"
	hideParentNameTag = "hideParentName"
	ifExistsRule      = "ifExists"
	omitemptyRule     = "omitempty"

	// TypeNames
	uuidType     = "uuid"
	uuidTypeName = "string UUID"

	// Delimiters
	fieldDelimiter = "."

	// Hints
	nestedPropsTag = "nestedProps"
)

var nestedPropsCompiler = regexp.MustCompile(`nestedProps=([a-zA-Z0-9|]+)`)

type Field interface {
	Name() string
	Value() interface{}
	TypeName() string
	ValidateIfExists() bool
	HideParentName() bool
	Omitempty() bool
	IsStruct() bool
	IsSlice() bool
	IsRequired() bool
	MustValidateType() bool
	Hints() []string

	IsValid() ([]rules.FieldError, bool)
	GenerateNestedFields() []Field
	GenerateRules() []rules.Rule
	ExtractValueFrom(data map[string]interface{}) interface{}

	SetName(value string)
	SetValue(value interface{})
	SetValidateIfExists(value bool)
}

type field struct {
	name               string
	value              interface{}
	typeName           string
	rules              []*rules.Rule
	validateIfExists   bool
	validationTagValue string
	jsonTagValue       string
	reflectType        reflect.StructField
	reflectValue       reflect.Value
}

func newField(fieldType reflect.StructField, fieldValue reflect.Value) Field {
	typeName := fieldType.Type.Kind().String()
	if fieldType.Type.Kind() == reflect.Slice {
		typeName = fmt.Sprintf("[]%s", fieldType.Type.Elem().Name())
	} else if strings.Contains(fieldValue.String(), uuidType) {
		typeName = uuidTypeName
	}
	name := strings.Split(fieldType.Tag.Get(jsonTag), ",")[0]
	validation := fieldType.Tag.Get(validationTag)
	validateIfExists := strings.Contains(validation, ifExistsRule) || !strings.Contains(validation, rules.REQUIRED)
	return &field{
		name:               name,
		typeName:           typeName,
		validationTagValue: validation,
		validateIfExists:   validateIfExists,
		jsonTagValue:       fieldType.Tag.Get(jsonTag),
		reflectType:        fieldType,
		reflectValue:       fieldValue,
	}
}

func (f *field) Name() string {
	return f.name
}

func (f *field) Value() interface{} {
	return f.value
}

func (f *field) TypeName() string {
	return f.typeName
}

func (f *field) ValidateIfExists() bool {
	return f.validateIfExists
}

func (f *field) SetName(value string) {
	f.name = value
}

func (f *field) SetValue(value interface{}) {
	f.value = value
}

func (f *field) SetValidateIfExists(value bool) {
	f.validateIfExists = value
}

func (f *field) IsValid() ([]rules.FieldError, bool) {
	var errs []rules.FieldError
	for _, rule := range f.GenerateRules() {
		if f.IsSlice() && !rule.IsSliceRule() {
			if f.value == nil || reflect.ValueOf(f.value).Len() == 0 {
				if rule.Type() == rules.REQUIRED {
					errs = append(errs, rule.GenerateError(f.Name()))
				}
				continue
			}
			for i := 0; i < reflect.ValueOf(f.value).Len(); i++ {
				element := reflect.ValueOf(f.value).Index(i).Interface()
				if !rule.IsValid(element) {
					errs = append(errs, rule.GenerateError(fmt.Sprintf("%s[%d]", f.Name(), i)))
				}
			}
		} else if !rule.IsValid(f.value) {
			errs = append(errs, rule.GenerateError(f.Name()))
			break
		}
	}
	return errs, len(errs) == 0
}

func (f *field) GenerateNestedFields() []Field {
	var nestedFieldInterfaceValue interface{}
	if f.IsStruct() {
		nestedFieldInterfaceValue = f.reflectValue.Interface()
	} else if f.IsSlice() {
		nestedFieldInterfaceValue = reflect.Zero(f.reflectValue.Type().Elem()).Interface()
	} else {
		return []Field{} // unknown type
	}
	nestedFields := buildValidators(nestedFieldInterfaceValue)
	var validNestedFieldNames []string
	if strings.Contains(f.validationTagValue, nestedPropsTag) {
		if ok := nestedPropsCompiler.Match([]byte(f.validationTagValue)); ok {
			validNestedFieldNames = nestedPropsCompiler.FindStringSubmatch(f.validationTagValue)[1:]
		}
	}
	validateIfHasValue := strings.Contains(f.reflectType.Tag.Get(validationTag), ifExistsRule)
	if len(validNestedFieldNames) > 0 {
		filteredFields := []Field{}
		for _, nestedField := range nestedFields {
			if slices.Contains(validNestedFieldNames, nestedField.Name()) {
				if validateIfHasValue && !nestedField.ValidateIfExists() {
					nestedField.SetValidateIfExists(validateIfHasValue)
				}
				filteredFields = append(filteredFields, nestedField)
				break
			}
		}
		nestedFields = filteredFields
	}
	for _, nestedField := range nestedFields {
		if f.HideParentName() {
			continue
		}
		nestedField.SetName(strings.ToLower(f.name) + fieldDelimiter + nestedField.Name())
	}
	return nestedFields
}

func (f *field) GenerateRules() []rules.Rule {
	var validators []rules.Rule
	if f.IsRequired() {
		validators = append(validators, rules.NewRequiredRule(f.typeName))
	}
	if f.IsRequired() || f.MustValidateType() {
		typeName := f.typeName
		if f.IsSlice() {
			typeName = f.reflectType.Type.Elem().Name()
		}
		if validator := rules.GetTypeValidator(typeName); validator != nil {
			validators = append(validators, validator)
		}
	}
	for _, hint := range f.Hints() {
		var validator rules.Rule = rules.GetRuleByHint(hint)
		if validator == nil {
			continue
		}
		validators = append(validators, validator)
	}
	return validators
}

func (f *field) ExtractValueFrom(data map[string]interface{}) interface{} {
	if !strings.Contains(f.name, fieldDelimiter) {
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			return data
		}
		return data[f.name]
	}
	var value interface{} = data
	for _, key := range strings.Split(f.name, fieldDelimiter) {
		if v, ok := value.(map[string]interface{}); ok {
			value = v[key]
		}
		if value == nil {
			return nil
		}
	}
	return value
}

func (f *field) HideParentName() bool {
	return f.reflectType.Tag.Get(hideParentNameTag) == "true"
}

func (f *field) Omitempty() bool {
	return strings.Contains(f.jsonTagValue, omitemptyRule)
}

func (f *field) IsStruct() bool {
	return f.reflectType.Type.Kind() == reflect.Struct
}

func (f *field) IsSlice() bool {
	return f.reflectType.Type.Kind() == reflect.Slice
}

func (f *field) IsRequired() bool {
	return strings.Contains(f.validationTagValue, rules.REQUIRED) && !f.Omitempty()
}

func (f *field) MustValidateType() bool {
	return strings.Contains(f.validationTagValue, rules.TYPE)
}

func (f *field) Hints() []string {
	return strings.Split(f.validationTagValue, ",")
}
