package configurator

import "encoding/json"

func DeserializeOptions[T any](data any) T {
	var out T
	raw, _ := json.Marshal(data)
	_ = json.Unmarshal(raw, &out)
	return out
}
