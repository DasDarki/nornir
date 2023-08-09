package annotations

import (
	"errors"
	"reflect"
	"strings"
)

func Decode(annotation *Annotation, target interface{}) error {
	if target == nil {
		return errors.New("target cannot be nil")
	}

	annotationDef, ok := target.(AnnotationDef)
	if !ok {
		return errors.New("target must be an AnnotationDef")
	}

	if annotationDef.GetName() != annotation.Name {
		return errors.New("annotation name does not match")
	}

	value := reflect.ValueOf(target).Elem()
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := value.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("annotation")
		if tag == "" {
			continue
		}

		props := make(map[string]string)
		parts := strings.Split(tag, ",")
		for _, part := range parts {
			keyValue := strings.SplitN(part, "=", 2)
			key, val := keyValue[0], keyValue[1]

			if key != "name" && key != "numeric" && key != "default" {
				return errors.New("invalid annotation tag")
			}

			props[key] = val
		}

		name := getPropString(props, "name")
		numeric := getPropString(props, "numeric")

		value, err := getDataValue(annotation, name, numeric)
		if err != nil {
			return err
		}

		if value == nil {
			defValue := getPropString(props, "default")
			if defValue == nil {
				return errors.New("missing annotation value")
			}

			if *defValue == "$empty" {
				value = getDefaultValueOfKind(fieldType.Type.Kind())
			} else {
				value = *defValue // currently only string supported
			}
		}

		field.Set(reflect.ValueOf(value))
	}

	if len(annotation.Data) > 0 {
		return errors.New("unused data in annotation")
	}

	return nil
}

func getDataValue(annotation *Annotation, name *string, numeric *string) (interface{}, error) {
	if name == nil && numeric == nil {
		return nil, errors.New("name and numeric cannot both be nil")
	}

	if name != nil {
		if value, ok := annotation.Data[*name]; ok {
			delete(annotation.Data, *name)
			return value, nil
		}
	}

	if numeric != nil {
		if value, ok := annotation.Data[*numeric]; ok {
			delete(annotation.Data, *numeric)
			return value, nil
		}
	}

	return nil, nil
}

func getPropString(props map[string]string, key string) *string {
	if value, ok := props[key]; ok {
		return &value
	}

	return nil
}

func getDefaultValueOfKind(kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Bool:
		return false
	case reflect.Int:
		return 0
	case reflect.Int8:
		return int8(0)
	case reflect.Int16:
		return int16(0)
	case reflect.Int32:
		return int32(0)
	case reflect.Int64:
		return int64(0)
	case reflect.Uint:
		return uint(0)
	case reflect.Uint8:
		return uint8(0)
	case reflect.Uint16:
		return uint16(0)
	case reflect.Uint32:
		return uint32(0)
	case reflect.Uint64:
		return uint64(0)
	case reflect.Uintptr:
		return uintptr(0)
	case reflect.Float32:
		return float32(0)
	case reflect.Float64:
		return float64(0)
	case reflect.Complex64:
		return complex64(0)
	case reflect.Complex128:
		return complex128(0)
	case reflect.Array:
		return [0]interface{}{}
	case reflect.Chan:
		return make(chan interface{})
	case reflect.Func:
		return func() {}
	case reflect.Interface:
		return interface{}(nil)
	case reflect.Map:
		return map[interface{}]interface{}{}
	case reflect.Ptr:
		return new(interface{})
	case reflect.Slice:
		return []interface{}{}
	case reflect.String:
		return ""
	case reflect.Struct:
		return struct{}{}
	default:
		return nil
	}
}
