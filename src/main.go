package main

import (
	"fmt"
	"listen"
)

func main() {

	fmt.Println("Starting application...")
	//go manager.start()

	//调用listen包，监听本地8001 端口
	err := listen.Listen("127.0.0.1:8001")
	if err != nil {
		fmt.Print(err)
	}

}
