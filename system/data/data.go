package data

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ToString(any interface{}) (string, error) {
	if any == nil || any == "" {
		return "", nil
	}

	switch v := any.(type) {
	case string:
		return any.(string), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	case fmt.Stringer:
		return v.String(), nil
	}

	return "", fmt.Errorf("Type cannot be converted to string: %s", reflect.TypeOf(any).String())
}

func ToBool(any interface{}) (bool, error) {
	if any == nil || any == "" {
		return false, nil
	}

	switch v := any.(type) {
	case bool:
		return v == true, nil
	case string:
		if v == "true" {
			return true, nil
		}
		if v == "false" {
			return false, nil
		}
	}

	return false, fmt.Errorf("Type cannot be converted to bool: %s", reflect.TypeOf(any).String())
}

func ToInt64(any interface{}) (int64, error) {
	if any == nil || any == "" {
		return 0, nil
	}

	switch v := any.(type) {
	case int:
		return int64(any.(int)), nil
	case int8:
		return int64(any.(int8)), nil
	case int16:
		return int64(any.(int16)), nil
	case int32:
		return int64(any.(int32)), nil
	case int64:
		return any.(int64), nil
	case string:
		return strconv.ParseInt(v, 10, 0)
	}

	return 0, fmt.Errorf("Type cannot be converted to int64: %s", reflect.TypeOf(any).String())
}

func LooseEqual(left, right interface{}) (bool, error) {
	if left == nil || left == "" {
		return right == nil || right == "", nil
	}
	if right == nil || right == "" {
		return left == nil || left == "", nil
	}

	switch l := left.(type) {
	case string:
		r, err := ToString(right)
		if err != nil {
			return false, err
		}
		return l == r, nil
	case int, int8, int16, int32, int64:
		li := int64(l.(int))
		r, err := ToInt64(right)
		if err != nil {
			return false, nil // Allow parse errors
		}
		return li == r, nil
	case bool:
		r, err := ToBool(right)
		if err != nil {
			return false, nil // Allow conversion error
		}
		return l == r, nil
	}

	return false, nil
}

func Compare(left, right interface{}) (int, error) {
	if left == nil || left == "" {
		if right == nil || right == "" {
			return 0, nil
		}
	}
	if right == nil || right == "" {
		if left == nil || left == "" {
			return 0, nil
		}
	}

	if left == nil || left == "" {
		return -1, nil
	}
	if right == nil || right == "" {
		return 1, nil
	}

	switch l := left.(type) {
	case string:
		r, err := ToString(right)
		if err != nil {
			return 0, err
		}
		return strings.Compare(l, r), nil
	case bool:
		r, err := ToBool(right)
		if err != nil {
			return 0, err
		}
		return CompareBool(l, r), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		li := int64(l.(int))
		r, err := ToInt64(right)
		if err != nil {
			return 0, nil // Allow parse errors
		}
		return int(li - r), nil
	}

	return 0, fmt.Errorf("Cannot compare types: %s, %s", reflect.TypeOf(left).String(), reflect.TypeOf(right).String())
}

func CompareBool(l, r bool) int {
	if l == r {
		return 0
	}
	if l {
		return 1
	}
	return -1
}
