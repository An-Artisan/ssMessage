package message

import "register"

type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content"`
}

func Send(message []byte, ignore *register.Client, manager *register.ClientManager) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}
