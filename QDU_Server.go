package main

// Importiert go-Bibliotheken (c# äquivalent: using System; ...)
import (
	// Standard Bibliotheken

	"database/sql"
	"log"
	"net"
	"os"

	"github.com/4kills/base64encoding"
	_ "github.com/go-sql-driver/mysql"
)

var config configuration
var db *sql.DB
var enc base64encoding.Encoder64

// Haupteinstiegspunkt des Programms beim Ausführen
func main() {

	// Startet thread (goroutine) die konstant den Benutzer-Input liest und auswertet
	readInput()

	// creates encoder for shorter links
	enc = base64encoding.New()

	// establishes connection with database
	initDB()
	defer db.Close()

	// Startet Web-Server, welcher konstant http-Anfragen verarbeitet
	// Bilder im Browser anzeigt, bei Aufrufen des links
	go webServer()

	log.Print("TCP-Server launched...\n\n")

	// Wartet auf TCP-Verbindungen durch den port, die Bilder auf den Server hochladen
	ln, err := net.Listen("tcp", config.PortTCP)
	if err != nil {
		log.Println("Fatal error:\n", err)
		os.Exit(1)
	}

	// Ermöglicht (beliebig) viele TCP-Verbindungen gleichzeitig
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// startet neue Goroutine der den verbundenen Benutzer bearbeitet
		go handleClient(conn)
	}
}

func initDB() {
	log.Print("DB-connection established... \n\n")
	var err error
	db, err = sql.Open("mysql", "4kills:4kills@/qdu") // später anpassen
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
}
