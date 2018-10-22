package protol

import (
	"encoding/json"
)

type Protocol struct {
	 Action string `json:"action"`
	 Body map[string]interface{} `json:"body"`
	 Headers map[string]string `json:"headers"`
}

// Convert a protocol to json string.
func (protocol *Protocol) ToString() (string, error) {
	bytes, err := json.Marshal(protocol)

	if err != nil {
		return "", err
	}  else {
		return string(bytes),nil
	}
}

// Create protocol from json string.
func FromJsonString(jsonString string) (*Protocol,error){

	protocol := &Protocol{}

	err := json.Unmarshal([]byte(jsonString), protocol)

	if err != nil {
		return nil, err
	}

	return protocol, nil
}