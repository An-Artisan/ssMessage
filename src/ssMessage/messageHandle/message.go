package messageHandle

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"time"
	"fmt"
)
// 定义消息最大长度
const (
	MaxMessageSize = 1024
)
// 定义消息体
type Message struct {
	Uid string `json:"uid"`
	MUid string `json:"muid"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content"`
}
// 定义内容解析格式
type Content struct {
	MessageContent string `json:"messageContent"`
	MUid   int `json:"muid"`
}

func (Manager *ClientManager) Send(message []byte, ignore *Client) {

	fmt.Println(Manager.Clients)
	// 广播除当前链接的用户链接信息
	for conn := range Manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

// 发送给自己
func (Manager *ClientManager) SendSelf(message []byte, conn *Client) {
	conn.Send <- message
}

// 发送数据给客户端
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
// 从客户端读取数据
func Read(conn *Client) {

	defer func() {
		Manager.Unregister <- conn
		conn.Socket.Close()
	}()
	conn.Socket.SetReadLimit(MaxMessageSize)
	for {
		_, message, err := conn.Socket.ReadMessage()

		if err != nil {
			Manager.Unregister <- conn
			conn.Socket.Close()
			break
		}
		// 如果是数据是心跳包则重置链接超时时间
		if string(message) == "OK"{
			conn.HeartBeat <- true
			continue
		}

		// 获取内容
		var content Content
		jsonErr  := json.Unmarshal([]byte(message),&content)

		// 消息格式不正确,则提示给客户端
		if jsonErr!=nil {
			Manager.MessageErr <- conn
			continue
		}
		// 判断用户数据库id是否存在
		if  content.MUid == 0{
			Manager.MessageErr <- conn
			timer := time.NewTimer(1 * time.Second)
			<-timer.C
			conn.Socket.Close()
		}
		// 收到信息广播给其他人
		jsonMessage, _ := json.Marshal(&Message{Uid: conn.Uid, Content: content.MessageContent})
		Manager.Broadcast <- jsonMessage
	}
}

