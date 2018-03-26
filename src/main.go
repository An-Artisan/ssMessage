package main

import (
	"fmt"
	"listen"
)

func main() {

	fmt.Println("Starting application...")
	//go manager.start()
	err:=listen.Listen("127.0.0.1:8001")
	if err != nil{
		fmt.Print(err)
	}

}