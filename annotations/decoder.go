package annotations

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Decode(annotation *Annotation, target interface{}, usage UsageKind) error {
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

	if !isUsageAllowed(annotationDef.GetUsages(), usage) {
		return errors.New("annotation usage is not allowed")
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
			key, val := strings.TrimSpace(keyValue[0]), strings.TrimSpace(keyValue[1])

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
				value, err = convertDefaultValueToKind(defValue, fieldType.Type.Kind())
				if err != nil {
					return err
				}
			}
		}

		field.Set(reflect.ValueOf(value))
	}

	if len(annotation.Data) > 0 {
		return fmt.Errorf("unused data in annotation: %v", annotation.Data)
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

func convertDefaultValueToKind(value *string, kind reflect.Kind) (interface{}, error) {
	switch kind {
	case reflect.String:
		return *value, nil
	case reflect.Bool:
		return *value == "true", nil
	case reflect.Int:
		return strconv.Atoi(*value)
	case reflect.Int8:
		return strconv.ParseInt(*value, 10, 8)
	case reflect.Int16:
		return strconv.ParseInt(*value, 10, 16)
	case reflect.Int32:
		return strconv.ParseInt(*value, 10, 32)
	case reflect.Int64:
		return strconv.ParseInt(*value, 10, 64)
	case reflect.Uint:
		return strconv.ParseUint(*value, 10, 0)
	case reflect.Uint8:
		return strconv.ParseUint(*value, 10, 8)
	case reflect.Uint16:
		return strconv.ParseUint(*value, 10, 16)
	case reflect.Uint32:
		return strconv.ParseUint(*value, 10, 32)
	case reflect.Uint64:
		return strconv.ParseUint(*value, 10, 64)
	case reflect.Uintptr:
		return strconv.ParseUint(*value, 10, 0)
	case reflect.Float32:
		return strconv.ParseFloat(*value, 32)
	case reflect.Float64:
		return strconv.ParseFloat(*value, 64)
	case reflect.Complex64:
		return strconv.ParseComplex(*value, 64)
	case reflect.Complex128:
		return strconv.ParseComplex(*value, 128)
	default:
		return nil, errors.New("unsupported kind")
	}
}

func isUsageAllowed(usages []UsageKind, usage UsageKind) bool {
	for _, u := range usages {
		if u == usage {
			return true
		}
	}

	return false
}
