package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    //"time"
    "github.com/gorilla/websocket"
    "github.com/satori/go.uuid"
)

type ClientManager struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

type Client struct {
    id     string
    socket *websocket.Conn
    send   chan []byte
}

type Message struct {
    Sender    string `json:"sender,omitempty"`
    Recipient string `json:"recipient,omitempty"`
    Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
    broadcast:  make(chan []byte),
    register:   make(chan *Client),
    unregister: make(chan *Client),
    clients:    make(map[*Client]bool),
}

func (Manager *ClientManager) start() {
    for {
        select {
        case conn := <-Manager.register:


            Manager.clients[conn] = true
            jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
            Manager.send(jsonMessage, conn)


        case conn := <-Manager.unregister:
            if _, ok := Manager.clients[conn]; ok {
                close(conn.send)
                delete(Manager.clients, conn)
                jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
                Manager.send(jsonMessage, conn)
            }

        case message := <-Manager.broadcast:
            for conn := range Manager.clients {
                select {
                case conn.send <- message:
                default:
                    close(conn.send)
                    delete(Manager.clients, conn)
                }
            }
        }
    }
}

func (Manager *ClientManager) send(message []byte, ignore *Client) {
    for conn := range Manager.clients {
        if conn != ignore {
            conn.send <- message
        }
    }
}

func (c *Client) read() {
    defer func() {
        Manager.unregister <- c
        c.socket.Close()
    }()
    //c.socket.SetReadDeadline(time.Now().Add(3*time.Second))
    for {
        _, message, err := c.socket.ReadMessage()
        if err != nil {
            Manager.unregister <- c
            c.socket.Close()
            break
        }
        jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
        Manager.broadcast <- jsonMessage
    }
}

func (c *Client) write() {
    defer func() {
        c.socket.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.socket.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            c.socket.WriteMessage(websocket.TextMessage, message)
        }
    }
}

func main() {
    fmt.Println("Starting application...")
    go Manager.start()
    http.HandleFunc("/ws", wsPage)
    http.ListenAndServe(":8001", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
    conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
    if error != nil {
        http.NotFound(res, req)
        return
    }

    sUuid,_:= uuid.NewV4()
    cUuid := uuid.NewV5(sUuid,"fdkljklfd").String()
    client := &Client{id: cUuid, socket: conn, send: make(chan []byte)}

    Manager.register <- client

    go client.read()
    go client.write()
}