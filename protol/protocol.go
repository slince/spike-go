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
func (protocol *Protocol) ToBytes() []byte {
	bytes, _ := json.Marshal(protocol)
	return bytes
}

// Convert a protocol to json string.
func (protocol *Protocol) ToString() string {
	return string(protocol.ToBytes())
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