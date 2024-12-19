package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Listening on port "6379")

	// Creating a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listening for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close() // Closing the connection

	for {
	buf := make([]byte, 1024)
	
	// Reading the data from the connection
	_, err = conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			break
		}
		fmt.Println("Error reading from client: ", err.Error())
		os.Exit(1)

	}

	// Ignore request if it is a PING and send back PONG
	conn.Write([]byte("+PONG\r\n"))


}