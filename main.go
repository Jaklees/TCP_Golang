package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("unable to start server: %s",err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Started server on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("unable to accept conect: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}