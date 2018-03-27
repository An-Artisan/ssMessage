package message

import "register"

type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content"`
}

<<<<<<< HEAD
func (manager *register.ClientManager) Send(message []byte, ignore *register.Client) {
=======
func Send(message []byte, ignore *register.Client, manager *register.ClientManager) {
>>>>>>> 44e1face9d923888421f132172ce11ce8365e7b0
	for conn := range manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}
