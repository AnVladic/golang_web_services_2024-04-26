package main

//
//import (
//	"errors"
//	"reflect"
//)
//
//func setFieldFromFloat(field reflect.Value, v reflect.Value) error {
//	switch field.Kind() {
//	case reflect.Int:
//		field.SetInt(int64(v.Float()))
//	case reflect.Slice:
//		if field.Type().Elem().Kind() == reflect.Int {
//			field.SetInt(int64(v.Float()))
//		} else {
//			field.Set(v)
//		}
//	default:
//		field.Set(v)
//	}
//	return nil
//}
//
//func reflectValue2Struct(
//	data interface{}, outField reflect.Value, valueOf func(i any) reflect.Value) error {
//	if !outField.CanSet() {
//		return nil
//	}
//
//	switch data.(type) {
//	case float64:
//		if err := setFieldFromFloat(outField, valueOf(data)); err != nil {
//			return err
//		}
//	case map[string]interface{}:
//		for k, v := range data.(map[string]interface{}) {
//			field := outField.FieldByName(k)
//			if err := reflectValue2Struct(v, field, reflect.ValueOf); err != nil {
//				return err
//			}
//		}
//	case []interface{}:
//		for _, v := range data.([]interface{}) {
//			appendValueOf := func(i any) reflect.Value {
//				return reflect.Append(outField, reflect.ValueOf(i))
//			}
//			if err := reflectValue2Struct(v, outField, appendValueOf); err != nil {
//				return err
//			}
//		}
//	default:
//		outField.Set(valueOf(data))
//	}
//	return nil
//}
//
//func i2s(data interface{}, out interface{}) error {
//	outValue := reflect.ValueOf(out).Elem()
//	if outValue.Kind() != reflect.Struct {
//		return errors.New("out is not a struct")
//	}
//	if err := reflectValue2Struct(data, outValue, reflect.ValueOf); err != nil {
//		return err
//	}
//	return nil
//}
