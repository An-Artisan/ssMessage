package connction


import (
	"net/http"
	"github.com/gorilla/websocket"
	"ssMessage/messageHandle"
	"time"
)

const(
	PongWait = time.Second * 5
)
func Listen(addr string) error {
	// 添加ws链接处理函数
	http.HandleFunc("/ws", WsHandle)
	//监听指定ip地址+端口
	err := http.ListenAndServe(addr, nil)
	return err
}

func WsHandle(res http.ResponseWriter, req *http.Request) {
	// 获得ws链接
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	// 出错返回404
	if err != nil {
		http.NotFound(res, req)
		return
	}
	// 注册链接,获取client
	client := messageHandle.SetUserInfo(conn)

	// 把cline写入到注册变量通道
	messageHandle.Manager.Register <- client


	// 连接进来开启一个协程读
	go messageHandle.Read(client)
	// 连接进来开启一个协程写
	go messageHandle.Write(client)

	go HeartBeat(client)
	client.HeartBeat <- true
}


// 心跳包检测
func HeartBeat(conn *messageHandle.Client) {
	for {
		select {
		case <-conn.HeartBeat:
			conn.Socket.SetReadDeadline(time.Now().Add(PongWait))
		case <-time.After(PongWait):
				messageHandle.Manager.Unregister <- conn
			return
		}
	}
}

