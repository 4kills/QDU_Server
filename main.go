package main

import (
	"fmt"
	"log"
	"net"

	"github.com/4kills/QDU_server/db"

	"github.com/4kills/base64encoding"
)

var config configuration
var enc base64encoding.Encoder64

func main() {
	// encoder for shorter links
	enc = base64encoding.New()

	// establishes connection with database
	if err := db.InitDB(); err != nil {
		log.Fatal(fmt.Errorf("DB connection error: ", err))
	}
	log.Println("Database connection established...")

	// starts web-server for http requests
	go webServer()
	log.Print("TCP-Server launched...\n\n")

	// listens for tcp connections through specified port and serves (pic upload, token)
	ln, err := net.Listen("tcp", config.PortTCP)
	if err != nil {
		log.Fatal("Fatal error:\n", err)
	}

	// allows for any number of concurrent tcp connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleClient(conn)
	}
}
