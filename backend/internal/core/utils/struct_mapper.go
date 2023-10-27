package utils

import (
	"encoding/json"
)

func MapJsonToStruct(message []byte, target interface{}) error {
	return json.Unmarshal(message, &target)
}
