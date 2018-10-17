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
	bytes, error := json.Marshal(protocol)

	if error != nil {
		return "", error
	}  else {
		return string(bytes),nil
	}
}

// Create protocol from json string.
func FromJsonString(jsonString string) (*Protocol,error){

	protocol := &Protocol{}

	error := json.Unmarshal([]byte(jsonString), protocol)

	if error != nil {
		return nil, error
	}

	return protocol, nil
}