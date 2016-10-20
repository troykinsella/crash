package util

import (
	"fmt"
)

type Values interface {
	Get(string) interface{}
	GetString(string) string
	AsMap() map[string]interface{}
}

type mapAdapter struct {
	vals map[string]interface{}
}

func (a *mapAdapter) Get(name string) interface{} {
	return a.vals[name]
}

func (a *mapAdapter) GetString(name string) string {
	return ToString(a.Get(name))
}

func (a *mapAdapter) AsMap() map[string]interface{} {
	return a.vals
}

func ToString(val interface{}) string {
	if val == nil || val == "" {
		return ""
	}

	switch val.(type) {
	case string:
		return val.(string)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case bool:
		return fmt.Sprintf("%t", val)
	}
	return ""
}

func AsValues(vals map[string]interface{}) Values {
	return &mapAdapter{
		vals: vals,
	}
}

func AsStringValues(vals map[string]string) Values {
	strs := make(map[string]interface{})
	for k, v := range vals {
		strs[k] = v
	}
	return &mapAdapter{
		vals: strs,
	}
}
