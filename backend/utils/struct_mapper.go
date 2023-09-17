package utils

import (
	"encoding/json"
)

func MapJsonToStruct(message []byte, target interface{}) error {
	err := json.Unmarshal(message, &target)
	return err
}
