package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/4kills/base64encoding"
	_ "github.com/go-sql-driver/mysql"
)

var config configuration
var db *sql.DB
var enc base64encoding.Encoder64

func main() {

	// starts goroutine listening for user input
	readInput()

	// encoder for shorter links
	enc = base64encoding.New()

	// establishes connection with database
	initDB()
	defer db.Close()

	// starts web-server for http requests
	go webServer()
	log.Print("TCP-Server launched...\n\n")

	// listens for tcp connections through specified port and serves (pic upload, token)
	ln, err := net.Listen("tcp", config.PortTCP)
	if err != nil {
		log.Println("Fatal error:\n", err)
		os.Exit(1)
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

func initDB() {
	log.Print("DB-connection established... \n\n")
	var err error

	if config.DBUser == "" {
		db, err = sql.Open("mysql", fmt.Sprintf("/%s", config.DBName))
	} else {
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", config.DBUser, config.DBPw, config.DBName))
	}

	if err != nil {
		log.Fatal("DB open error:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("DB connection error: ping unsuccessful:", err)
	}
}
