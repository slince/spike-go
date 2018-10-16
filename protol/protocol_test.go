package protol

import "testing"

func TestFromJsonString(t *testing.T) {
	json := `{action: "register", }`
	protocol := FromJsonString(json)
}