package ast

import (
	"go/types"
	"reflect"
	"net/http"
)

func extractValue(operand interface{}, key interface{}) (interface{}, error) {
	switch o := operand.(type) {
	case http.Header:
		return o.Get(key.(string)), nil
	case types.Map:
		v := reflect.ValueOf(operand)
		kv := reflect.ValueOf(key)
		return v.MapIndex(kv), nil
	}

	return nil, nil
}
