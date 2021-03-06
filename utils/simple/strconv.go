package simple

import (
	"reflect"
)

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice:
		//use of IsNil method
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func InString(hay *[]string, needle string) bool {
	for _, x := range *hay {
		if x == needle {
			return true
		}
	}
	return false
}

func IsArray(a interface{}) bool {
	rt := reflect.TypeOf(a)
	if rt == nil {
		return false
	}
	return rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array
}

func ToString(x interface{}) (string, bool) {
	switch t := x.(type) {
	case string:
		return t, true
	case *string:
		return *t, true
	case []byte:
		return string(t), true
	default:
		return "", false
	}
}
func ToFloat(x interface{}) (float64, bool) {
	switch t := x.(type) {
	case bool:
		if t {
			return 1, true
		} else {
			return 0, true
		}
	case float64:
		return t, true
	case int:
		return float64(t), true
	case uint:
		return float64(t), true

	case int8:
		return float64(t), true
	case int16:
		return float64(t), true
	case uint8:
		return float64(t), true
	case uint16:
		return float64(t), true

	case uint32:
		return float64(t), true
	case int32:
		return float64(t), true
	case uint64:
		return float64(t), true
	case int64:
		return float64(t), true
	case *float64:
		return *t, true
	case *int:
		return float64(*t), true
	case *uint:
		return float64(*t), true

	case *int8:
		return float64(*t), true
	case *int16:
		return float64(*t), true
	case *uint8:
		return float64(*t), true
	case *uint16:
		return float64(*t), true

	case *uint32:
		return float64(*t), true
	case *int32:
		return float64(*t), true
	case *uint64:
		return float64(*t), true
	case *int64:
		return float64(*t), true

	default:
		return 0, false
	}
}

// Does x resemble a int?
func IsInt(x interface{}) bool {
	switch x.(type) {
	case bool, int, int8, int16, int32, int64,
		uint8, uint16, uint32, uint64:
		return true
	}

	return false
}

func ToInt64(x interface{}) (int64, bool) {
	switch t := x.(type) {
	case bool:
		if t {
			return 1, true
		} else {
			return 0, true
		}
	case int:
		return int64(t), true
	case uint8:
		return int64(t), true
	case int8:
		return int64(t), true
	case uint16:
		return int64(t), true
	case int16:
		return int64(t), true
	case uint32:
		return int64(t), true
	case int32:
		return int64(t), true
	case uint64:
		return int64(t), true
	case int64:
		return t, true
	case float64:
		return int64(t), true
	case *int:
		return int64(*t), true
	case *uint8:
		return int64(*t), true
	case *int8:
		return int64(*t), true
	case *uint16:
		return int64(*t), true
	case *int16:
		return int64(*t), true
	case *uint32:
		return int64(*t), true
	case *int32:
		return int64(*t), true
	case *uint64:
		return int64(*t), true
	case *int64:
		return int64(*t), true
	case *float64:
		return int64(*t), true
	default:
		return 0, false
	}
}
