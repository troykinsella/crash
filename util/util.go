package util

func KeysForMapStringInterface(obj map[string]interface{}) []string {
	result := make([]string, len(obj))
	i := 0
	for k := range obj {
		result[i] = k
		i++
	}
	return result
}

func PutAllForMapStringString(dest map[string]string, source map[string]string) {
	for k, v := range source {
		dest[k] = v
	}
}

func PutAllForMapStringInterface(dest map[string]interface{}, source map[string]interface{}) {
	for k, v := range source {
		dest[k] = v
	}
}
