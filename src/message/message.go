package message

import "register"


type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content"`
}

func (manager *register.ClientManager) Send(message []byte, ignore *register.Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}