package validator

import (
	"encoding/json"
	"github.com/wallrony/go-validator/rules"
	"reflect"
)

func ValidateDTOPartially[T interface{}](data interface{}) (*T, ValidationError) {
	return buildGenericInstance[T](data), validate[T](data)
}

func ValidateDTO[T interface{}](data interface{}) (*T, ValidationError) {
	err := validate[T](data)
	if err != nil {
		return nil, err
	}
	return buildGenericInstance[T](data), nil
}

func buildValidators(instance interface{}) []Field {
	var reflection = reflect.ValueOf(instance)
	if reflection.Kind() == reflect.Ptr {
		reflection = reflection.Elem()
	}
	var fields []Field
	for i := 0; i < reflection.NumField(); i++ {
		fieldValue := reflection.Field(i)
		fieldType := reflection.Type().Field(i)
		field := newField(fieldType, fieldValue)
		fields = append(fields, field)
		if field.IsStruct() || field.IsSlice() {
			fields = append(fields, field.GenerateNestedFields()...)
		}
	}
	return fields
}

func tryValidators(data interface{}, fields []Field) []rules.FieldError {
	var errs []rules.FieldError
	var formattedData map[string]interface{} = formatJSONData(data)
	for _, field := range fields {
		if field.IsStruct() {
			continue
		}
		value := field.ExtractValueFrom(formattedData)
		canJump := field.ValidateIfExists() && value == nil
		if value != nil && field.TypeName() == reflect.Slice.String() {
			canJump = field.ValidateIfExists() && reflect.ValueOf(value).Len() == 0
		}
		if canJump {
			continue
		}
		field.SetValue(value)
		if fieldErrs, ok := field.IsValid(); !ok {
			errs = append(errs, fieldErrs...)
		}
	}
	return errs
}

func validate[T interface{}](data interface{}) ValidationError {
	var it T
	var validators []Field = buildValidators(it)
	var fieldsErrors = tryValidators(data, validators)
	if len(fieldsErrors) == 0 {
		return nil
	}
	return newValidationError(fieldsErrors)
}

func buildGenericInstance[T interface{}](data interface{}) *T {
	dataStr, _ := json.Marshal(formatJSONData(data))
	var instance T
	json.Unmarshal(dataStr, &instance)
	return &instance
}
