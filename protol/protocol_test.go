package protol

import (
	"testing"
)

func TestFromJsonString(t *testing.T) {
	json := `{"action": "register", "headers": {"status": "200"}, "body": {"error": ""}}`
	protocol, err := CreateFromJson(json)

	if err != nil {
		t.Errorf("error parse json")
	}

	if protocol.Action != "register" {
		t.Errorf("error parse action")
	}

	if header, ok := protocol.Headers["status"]; !ok || header != "200"{
		t.Errorf("error parse header")
	}
	if body, ok := protocol.Body["error"]; !ok || body != ""{
		t.Errorf("error parse body")
	}
}

func TestToString(t *testing.T) {
	protocol := Protocol{
		Action: "register",
		Headers: map[string]string{
			"status": "200",
		},
		Body:map[string]interface{}{
			"error": "",
		},
	}

	expectedJson := `{"action":"register","body":{"error":""},"headers":{"status":"200"}}`

	if json := protocol.ToString(); json != expectedJson {
		t.Errorf("error tostring")
	}
}