package messageHandle

import (
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"encoding/json"

)

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
}

type Client struct {
	Uid    string
	MUid   int
	Socket *websocket.Conn
	Send   chan []byte
}

func GetUid() string {
	sUuid, _ := uuid.NewV4();
	cUuid := uuid.NewV5(sUuid, "fdkljklfd").String()
	return cUuid
}

func SetUserInfo(conn *websocket.Conn) *Client {
	client := &Client{Uid: GetUid(), MUid: 0, Socket: conn, Send: make(chan []byte)}
	return client
}
func (manager *ClientManager) Start() {
	// 监听用户动作
	for {
		select {
		// 注册时，发送登录信息广播给其他用户
		case conn := <-manager.Register:
			// 给Clinets 值赋值为 true
			manager.Clients[conn] = true
			// 组装json数据
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			// 开始发送数据

			manager.Send(jsonMessage, conn)

			//	注销或断开链接,发送退出信息广播给其他用户
		case conn := <-manager.Unregister:
			// 判断conn信息是否存在
			if _, ok := manager.Clients[conn]; ok {
				//关闭Send通道
				close(conn.Send)
				//删除该用户链接
				delete(Manager.Clients, conn)

				// 组装json数据
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				// 发送广播数据
				manager.Send(jsonMessage, conn)
			}

			//	 接受信息
		case messageContent := <-manager.Broadcast:
			for conn := range manager.Clients {
				select {
				case conn.Send <- messageContent:
				default:
					close(conn.Send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}
