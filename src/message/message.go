package message

import "register"


type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content"`
}
var C  register.ClientManager
func (manager *C) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}