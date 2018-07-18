package main

import (
	"fmt"
	"reflect"
)

const (
	validationTagName      = "required"
	validationErrorMessage = "field %s must be set"
)

func invalidField(fieldName string) error {
	return fmt.Errorf(validationErrorMessage, fieldName)
}

// Validate that option structs have required options
func Validate(options interface{}) error {
	val := reflect.ValueOf(options)

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("Invalid input type %s, expected Struct", val.Kind())
	}

	structTyp := reflect.TypeOf(options)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := structTyp.Field(i)
		tag := fieldType.Tag.Get(validationTagName)

		if tag == "true" {
			switch field.Kind() {
			case reflect.Int:
				if field.Int() == 0 {
					return invalidField(fieldType.Name)
				}
			case reflect.String:
				fallthrough
			case reflect.Slice:
				fallthrough
			case reflect.Array:
				if field.Len() == 0 {
					return invalidField(fieldType.Name)
				}
			case reflect.Struct:
				if err := Validate(field.Interface()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
