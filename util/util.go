package util

import "time"

func KeysForMapStringInterface(obj map[string]interface{}) []string {
	result := make([]string, len(obj))
	i := 0
	for k := range obj {
		result[i] = k
		i++
	}
	return result
}

func PutAllForMapStringInterface(dest map[string]interface{}, source map[string]interface{}) {
	for k, v := range source {
		dest[k] = v
	}
}

func Timeout(d time.Duration) chan bool {
	if d <= 0 {
		return nil
	}

	t := make(chan bool, 1)
	go func() {
		time.Sleep(d)
		t <- true
	}()
	return t
}
