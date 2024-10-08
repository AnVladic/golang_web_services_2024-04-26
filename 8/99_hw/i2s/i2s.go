package main

import (
	"errors"
	"reflect"
)

func reflectValue2Struct(value interface{}, outValue reflect.Value) error {
	switch outValue.Kind() {
	case reflect.Int:
		val, ok := value.(float64)
		if !ok {
			return errors.New("value is not float64")
		}
		outValue.SetInt(int64(val))
	case reflect.Bool:
		val, ok := value.(bool)
		if !ok {
			return errors.New("value is not bool")
		}
		outValue.SetBool(val)
	case reflect.String:
		val, ok := value.(string)
		if !ok {
			return errors.New("value is not string")
		}
		outValue.SetString(val)
	case reflect.Struct:
		t := outValue.Type()
		dataMap, ok := value.(map[string]interface{})
		if !ok {
			return errors.New("data is not a map")
		}

		for i := range outValue.NumField() {
			field := outValue.Field(i)
			name := t.Field(i).Name
			value := dataMap[name]
			if err := reflectValue2Struct(value, field); err != nil {
				return err
			}
		}
	case reflect.Slice:
		dataSlice, ok := value.([]interface{})
		if !ok {
			return errors.New("data is not a slice")
		}
		sliceType := outValue.Type().Elem()
		newSlice := reflect.MakeSlice(outValue.Type(), 0, len(dataSlice))

		for _, data := range dataSlice {
			newElem := reflect.New(sliceType).Elem()
			if err := reflectValue2Struct(data, newElem); err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, newElem)
		}
		outValue.Set(newSlice)
	default:
		outValue.Set(reflect.ValueOf(value))
	}
	return nil
}

func i2s(data interface{}, out interface{}) error {
	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr {
		return errors.New("outValue is not a pointer")
	}
	outValue = outValue.Elem()

	if err := reflectValue2Struct(data, outValue); err != nil {
		return err
	}
	return nil
}
