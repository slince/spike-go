package client

type Client struct {
	// Client address
	Address string
	// Server Address
	ServerAddress string
	// The folder to listen
	LocalPath string
}

// Run client
func (client *Client) Run() {

}

// Sync files from server.
func (client *Client) Sync() {

}

// Send file to serer.
func (client *Client) sendFile() error{

}