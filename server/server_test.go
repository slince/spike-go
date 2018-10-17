package server

import "testing"

func TestCreateServer(t *testing.T) {
	server := CreateServer("127.0.0.1:9090")

	if server == nil {
		t.Errorf("Fail to create a server")
	}
}
