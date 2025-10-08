package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"redis-clone/config"
	"redis-clone/core"
)

func RunSyncTcpServer() {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	fmt.Printf("Starting Server at %s", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		for {
			commands, err := readCommands(conn)
			// conn.Read()
			if err != nil {
				conn.Close()
				if err == io.EOF {
					break
				}
				log.Println("client_disconnnected")
			}
			fmt.Println("OUTPUT", commands)
			//Eval the command
			//write the output
			core.EvalAndRespond(commands, conn)
		}

	}
}
