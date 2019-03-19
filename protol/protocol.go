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
func CreateFromJson(jsonString string) (*Protocol,error){
	return CreateFromBytes([]byte(jsonString))
}

// Create protocol from bytes buffer.
func CreateFromBytes(buffer []byte) (*Protocol,error){
	protocol := &Protocol{}
	err := json.Unmarshal(buffer, protocol)

	if err != nil {
		return nil, err
	}
	return protocol, nil
}