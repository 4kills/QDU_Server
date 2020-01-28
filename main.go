package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/4kills/base64encoding"

	"github.com/4kills/qdu_server/db"
	"github.com/4kills/qdu_server/web"
)

var enc base64encoding.Encoder64

func main() {
	enc = base64encoding.New()

	// establishes connection with database
	if err := db.InitDB(); err != nil {
		log.Fatal(fmt.Errorf("DB connection error: %s", err))
	}
	log.Println("Database connection established. . .")

	// starts web-server for http requests
	go web.Server()
	log.Println("Web-Server launched. . .")

	// listens for tcp connections through specified port and serves (pic upload, token)
	ln, err := net.Listen("tcp", os.Getenv("PORT_TCP"))
	if err != nil {
		log.Fatal("Fatal error: ", err)
	}
	log.Println("TCP-Server launched: Listening for incomming screenshots. . .")

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
