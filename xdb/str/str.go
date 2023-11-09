package str

import "encoding/json"

func AnyToJson(obj interface{}) string {
	barr, err := json.Marshal(obj)
	if err != nil {
		return "failed to Marshal: " + err.Error()
	}
	return string(barr)
}
