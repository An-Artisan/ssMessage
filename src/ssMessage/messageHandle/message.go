package messageHandle

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"fmt"
)


// 定义消息体
type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content"`
}

func (manager *ClientManager) Send(message []byte, ignore *Client) {

	for conn ,_:= range manager.Clients {
		fmt.Println(conn.Uid )
		fmt.Println("===" )
	}
	// 广播除当前链接的用户链接信息
	for conn ,_:= range manager.Clients {
		if conn != ignore {
			conn.Send <- message

		}
	}
}

func Write(conn *Client) {
	// 程序结束后关闭链接
	defer func() {
		conn.Socket.Close()
	}()
	// 轮询监听
	for {
		select {
		// 如果有通道有信息就发送,通道关闭就发送关闭信息 (这里 <-conn.Send第二个结果是一个bool型,false代表通道被关闭)
		case message, ok := <-conn.Send:
			if !ok{
				conn.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			conn.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func  Read(conn *Client) {
	defer func() {
		Manager.Register <- conn
		conn.Socket.Close()
	}()
	//c.socket.SetReadDeadline(time.Now().Add(3*time.Second))
	for {
		_, message, err := conn.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- conn
			conn.Socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: conn.Uid, Content: string(message)})
		Manager.Broadcast <- jsonMessage
	}
}


