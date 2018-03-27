package main

import (
	"fmt"

	"ssMessage/messageHandle"
	"ssMessage/connction"
)

func main() {

	fmt.Println("Starting application...")
	//go Manager.start()
	go messageHandle.Manager.Start()

	//调用listen包，监听本地8001 端口
	err := connction.Listen("127.0.0.1:8001")

	if err != nil {
		fmt.Print(err)
	}

}
