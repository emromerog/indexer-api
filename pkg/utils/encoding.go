package utils

import (
	"encoding/json"
	"fmt"
)

func ConvertToJson(body interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error converting TO json: %v", err)
	}
	return jsonData, nil
}

func ConvertFromJSON(jsonData []byte, v interface{}) error {
	err := json.Unmarshal(jsonData, v)
	if err != nil {
		return fmt.Errorf("error converting FROM json: %v", err)
	}
	return nil
}
