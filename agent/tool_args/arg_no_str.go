package toolargs

import "encoding/json"

func checkArgNotString(args string) string {
	var a map[string]any
	err := json.Unmarshal([]byte(args), &a)
	if err != nil {
		return ""
	}

	val, ok := a["__arg1"]
	if !ok {
		return ""
	}

	bytes, _ := json.Marshal(val)
	return string(bytes)
}
