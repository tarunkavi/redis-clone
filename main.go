package main

import (
	"flag"
	"fmt"
	"log"
	"redis-clone/config"
	"redis-clone/core"
	"redis-clone/server"
)

func init() {
	// Set log flags to include file and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for redis server")
	flag.IntVar(&config.Port, "port", 7379, "host for redis server")
	flag.Parse()
}

func main() {
	setupFlags()
	fmt.Println("Starting db....")
	core.InitAof(config.AOFFile)
	err := server.RunAsyncTcpServer()
	fmt.Println(err)
	// out, _ := core.Decode([]byte("*-1\r\n"))
	// fmt.Println(out)

}

//
