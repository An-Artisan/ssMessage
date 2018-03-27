package connction


import (
	"net/http"
	"github.com/gorilla/websocket"
	"ssMessage/messageHandle"
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

	go messageHandle.Read(client)
	go messageHandle.Write(client)
}

