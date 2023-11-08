package ptr

import "encoding/json"

func Any[T any](obj T) *T {
	return &obj
}

func AnyToJson(obj interface{}) string {
	barr, err := json.Marshal(obj)
	if err != nil {
		return "failed to Marshal: " + err.Error()
	}
	return string(barr)
}
