package validate

import (
	"errors"
	"fmt"
	"homework/internal/domain"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrNotStruct                   = errors.New("wrong argument given, should be a struct")
	ErrInvalidValidatorSyntax      = errors.New("invalid validator syntax")
	ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")
	ErrLenValidationFailed         = errors.New("len validation failed")
	ErrInValidationFailed          = errors.New("in validation failed")
	ErrMaxValidationFailed         = errors.New("max validation failed")
	ErrMinValidationFailed         = errors.New("min validation failed")
)

type ValidationError struct {
	field string
	err   error
}

func NewValidationError(err error, field string) error {
	return &ValidationError{
		field: field,
		err:   err,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.field, e.err)
}

func (e *ValidationError) Unwrap() error {
	return e.err
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var errs []error

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		if !fieldValue.CanInterface() {
			errs = append(errs, NewValidationError(ErrValidateForUnexportedFields, field.Name))
			return errors.Join(errs...)
		}
		errs = append(errs, validateField(field.Name, fieldValue.Interface(), tag)...)
	}
	return errors.Join(errs...)
}

func validateField(fieldName string, fieldValue interface{}, tag string) []error {
	parts := strings.Split(tag, ":")
	if len(parts) != 2 {
		return []error{NewValidationError(ErrInvalidValidatorSyntax, fieldName)}
	}

	validator := parts[0]
	value := parts[1]

	switch fv := fieldValue.(type) {
	case []int:
		errs := make([]error, 0)
		for _, v := range fv {
			errs = append(errs, validateFieldValue(validator, value, fieldName, v))
		}
		return errs
	case []string:
		errs := make([]error, 0)
		for _, v := range fv {
			errs = append(errs, validateFieldValue(validator, value, fieldName, v))
		}
		return errs
	default:
		return []error{validateFieldValue(validator, value, fieldName, fv)}
	}
}

func validateFieldValue(validator, value, fieldName string, fieldValue interface{}) error {
	switch validator {
	case "len":
		return validLen(fieldName, value, fieldValue)
	case "in":
		return validIn(fieldName, value, fieldValue)
	case "min":
		return validMin(fieldName, value, fieldValue)
	case "max":
		return validMax(fieldName, value, fieldValue)
	}
	return nil
}

func validLen(fieldName, value string, fieldValue interface{}) error {
	length, err := strconv.Atoi(value)
	if err != nil {
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}
	if length < 0 {
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}

	if reflect.TypeOf(fieldValue).Kind() != reflect.String {
		return NewValidationError(ErrLenValidationFailed, fieldName)
	}

	if len(fieldValue.(string)) != length {
		return NewValidationError(ErrLenValidationFailed, fieldName)
	}
	return nil
}

func validIn(fieldName, value string, fieldValue interface{}) error {
	if len(value) == 0 {
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}

	validValues := strings.Split(value, ",")
	valid := false

	if validValues[0] == "SensorType" {
		if fmt.Sprintf("%v", fieldValue) == string(domain.SensorTypeContactClosure) || fmt.Sprintf("%v", fieldValue) == string(domain.SensorTypeADC) {
			valid = true
		}
	} else {
		for _, validValue := range validValues {
			if fmt.Sprintf("%v", fieldValue) == validValue {
				valid = true
				break
			}
		}
	}

	if !valid {
		return NewValidationError(ErrInValidationFailed, fieldName)
	}
	return nil
}

func validMin(fieldName, value string, fieldValue interface{}) error {
	min, err := strconv.Atoi(value)
	if err != nil {
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}

	if reflect.TypeOf(fieldValue).Kind() == reflect.Ptr {
		fieldValue = reflect.ValueOf(fieldValue).Elem().Interface()
	}

	switch fieldValue := fieldValue.(type) {
	case string:
		if len(fieldValue) < min {
			return NewValidationError(ErrMinValidationFailed, fieldName)
		}
	case int:
		if fieldValue < min {
			return NewValidationError(ErrMinValidationFailed, fieldName)
		}
	case int64:
		if fieldValue < int64(min) {
			return NewValidationError(ErrMaxValidationFailed, fieldName)
		}
	default:
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}
	return nil
}

func validMax(fieldName, value string, fieldValue interface{}) error {
	max, err := strconv.Atoi(value)
	if err != nil {
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}

	switch reflect.TypeOf(fieldValue).Kind() {
	case reflect.String:
		if len(fieldValue.(string)) > max {
			return NewValidationError(ErrMaxValidationFailed, fieldName)
		}
	case reflect.Int:
		if fieldValue.(int) > max {
			return NewValidationError(ErrMaxValidationFailed, fieldName)
		}
	case reflect.Int64:
		if fieldValue.(int64) > int64(max) {
			return NewValidationError(ErrMaxValidationFailed, fieldName)
		}
	default:
		return NewValidationError(ErrInvalidValidatorSyntax, fieldName)
	}
	return nil
}
