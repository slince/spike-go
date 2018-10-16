package protol

import "encoding/json"

type Protocol struct {
	 Action string
	 Body map[string]interface{}
	 Headers map[string]interface{}
}

// Make protocol to json.
func (protocol *Protocol) ToString() (string, error) {
	bytes, error := json.Marshal(protocol)

	if error != nil {
		return "", error
	}  else {
		return string(bytes),nil
	}
}

// Create protocol from json
func FromJsonString(jsonString string) Protocol{

}