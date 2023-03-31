package main

import (
	"fmt"
	"reflect"

	"github.com/kardianos/service"
)

func setArguments(s service.Service, args []string) error {
	sValue := reflect.ValueOf(s)
	if sValue.Kind() == reflect.Ptr {
		sValue = sValue.Elem()
	}

	if sValue.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct: %v", sValue.Type())
	}

	configField := sValue.FieldByName("Config")
	if !configField.IsValid() {
		return fmt.Errorf("config field not found")
	}

	if configField.Kind() != reflect.Ptr || configField.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config field is not a pointer to a struct")
	}

	argumentsField := configField.Elem().FieldByName("Arguments")
	if !argumentsField.IsValid() {
		return fmt.Errorf("arguments field not found")
	}

	if argumentsField.Kind() != reflect.Slice {
		return fmt.Errorf("arguments field is not a slice")
	}

	argumentsField.Set(reflect.ValueOf(args))
	return nil
}
