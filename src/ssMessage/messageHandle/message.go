package messageHandle

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"time"
	"fmt"
)

// 定义消息体
type Message struct {
	Uid string `json:"uid"`
	MUid string `json:"muid"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content"`
}

type Content struct {
	MessageContent string `json:"messageContent"`
	HeartBeat   bool `json:"heartBeat"`
}

func (Manager *ClientManager) Send(message []byte, ignore *Client) {

	// 广播除当前链接的用户链接信息
	for conn := range Manager.Clients {
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
			if !ok {
				conn.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			conn.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func Read(conn *Client) {
	defer func() {
		Manager.Unregister <- conn
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
		data := make(chan byte)
		go getMessage([]byte(message), data)
		go heartBeat(conn,data,4)
		var content Content
		jsonErr  := json.Unmarshal([]byte(message),&content)
		if jsonErr!=nil {
			Manager.MessageErr <- conn
		}

		// 收到信息广播给其他人
		jsonMessage, _ := json.Marshal(&Message{Uid: conn.Uid, Content: string(message)})
		Manager.Broadcast <- jsonMessage
	}
}
func getMessage(message []byte, data chan byte) {
	for _, v := range message {
		data <- v
	}
	close(data)
}
func heartBeat(conn *Client, data chan byte, timeout int) {
	select {
	case <-data:
		//beego.Trace(conn.RemoteAddr().String(), string(fk), ": ", "Heart beating ")
		conn.Socket.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	case <-time.After(time.Second * 5):
		Manager.HeartBeat <- conn
		conn.Socket.Close()
	}
}
