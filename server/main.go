package main

import (
	"bufio"
	"log"
	"net"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	openConnections = make(map[net.Conn]bool)
	newConnection   = make(chan net.Conn)
	deadConnection  = make(chan net.Conn)
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	logFatal(err)

	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			logFatal(err)

			openConnections[conn] = true
			newConnection <- conn
		}
	}()

	for {
		select {
		case conn := <-newConnection:
			//Invoke broadcase message (broadcases to the other connections)
			go broadcastMessage(conn)

		case conn := <-deadConnection:
			//remove/delete the connection
			for item := range openConnections {
				if item == conn {
					break
				}
			}

			delete(openConnections, conn)
		}
	}
}

func broadcastMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		//loop through all the open connections
		//and send messages to these connections
		//except the connection that sent the message

		for item := range openConnections {
			if item != conn {
				item.Write([]byte(message))
			}
		}
	}

	deadConnection <- conn
}
