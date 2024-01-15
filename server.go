package main

import (
	"fmt"
	"net"
)

func ErrorManager(err error) {
	if err != nil {
		panic(err)
		// fmt.Println("Error")
	}
}

const (
	IP   = "127.0.0.1"
	PORT = "1212"
)

func mains() {
	fmt.Println("Server is running...")

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	ErrorManager(err)

	conn, err := ln.Accept()
	ErrorManager(err)

	fmt.Println("A client is connected by", conn.RemoteAddr())

	for {
		buffer := make([]byte, 4096)
		length, err := conn.Read(buffer)
		msg := string(buffer[:length])
		if err != nil {
			fmt.Println("Client leaved chat...")
		}
		fmt.Print("Client:", msg)

		conn.Write([]byte(msg))
	}
}
