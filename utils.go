package main

import (
	"errors"
	"fmt"
	"reflect"
)

func SetField(obj interface{}, name string, value interface{}) error {
	fmt.Printf("Obj = %v, Name = %v, Value = %v", &obj, name, value)
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match object field type")
	}

	structFieldValue.Set(val)
	return nil
}

func FillStruct(s interface{}, m map[string]interface{}) error {
	fmt.Printf("%v", m)
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
