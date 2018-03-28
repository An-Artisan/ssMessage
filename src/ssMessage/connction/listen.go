package connction


import (
	"net/http"
	"github.com/gorilla/websocket"
	"ssMessage/messageHandle"
	"time"
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

	//client.Socket.WriteMessage(websocket.TextMessage, []byte ("HelloWorld"))
	//fmt.Print(client)

	// 连接进来开启一个协程读
	go messageHandle.Read(client)
	// 连接进来开启一个协程写
	go messageHandle.Write(client)
	//data := make(chan  int)
	//go getMessage(messageHandle.HeartMessage,data)

	go HeartBeat(client,4)
	messageHandle.HeartMessage <- 1

}

func HeartBeat(conn *messageHandle.Client, imeout int) {

	for {
		select {
		case <-messageHandle.HeartMessage:
			//beego.Trace(conn.RemoteAddr().String(), string(fk), ": ", "Heart beating ")
			conn.Socket.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))
		case <-time.After(time.Second * 5):
			messageHandle.Manager.HeartBeat <- conn
			return

		}

	}

}

