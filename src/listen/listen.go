package listen

import (
	"net/http"
	"github.com/gorilla/websocket"
	"register"
	"fmt"
)

func Listen(addr string) error{

	http.HandleFunc("/ws", WsHandle)
	err := http.ListenAndServe(addr, nil)

	return err
}


func WsHandle(res http.ResponseWriter, req *http.Request) {

	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}

	client := register.SetUserInfo(conn)

	register.Manager.Register <- client
	register.Manager.Start()
	//client.Socket.WriteMessage(websocket.TextMessage, []byte ("HelloWorld"))
	//fmt.Print(client)
	//manager.register <- client

	//go client.read()
	//go client.write()
}

