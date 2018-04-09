//---------------------------------------------------------
// Go/Golang wurde von mir als Server-Sided-Scripting-Language
// gewählt, weil es mit jedem(!) geläufigen Betriebssystem kompatibel
// ist. Während c# mit dotnet core 2.x ebenfalls Unterstützung für
// Linux-arm und Linux-x64 anbietet, fehlt nach wie vor
// Linux-x86 Unterstützung.
// Darüber hinaus ist Go extrem effizient, bezüglich Garbage-
// Collection und Networking. Go ist außerdem eine strukturierte,
// imperative Programmiersprache mit großartiger Implementierung
// von Nebenläufigkeit (Co-Routinen) als go's Goroutines.
//---------------------------------------------------------

// "Namespace" in go; muss Funktion func main(){} enthalten
// (Haupteinstiegspunkt des Programms)
package main

// Importiert go-Bibliotheken (c# äquivalent: using System; ...)
import (
	// Standard Bibliotheken
	"fmt"
	"net"
	"os"
	"time"
)

// Attribute; Gültig über alle Funktionen
// Asynchron setzbar von jeder Goroutine
// (c#'s threads, co-routinen, keine sub-routinen)
var domain string
var directory string
var directoryWeb string
var port string
var portWeb string

// Haupteinstiegspunkt des Programms beim Ausführen
func main() {
	// Channel: werden benutzt um verschieden Goroutines zu synchronisieren
	portChan := make(chan string)
	portChanWeb := make(chan string)
	dirChan := make(chan string)

	// Startet thread (goroutine) die konstant den Benutzer-Input liest und auswertet
	go readInput(portChan, portChanWeb, dirChan)
	// Startet Web-Server, welcher konstant http-Anfragen verarbeitet
	// Bilder im Browser anzeigt, bei Aufrufen des links
	go webServer(portChanWeb, dirChan)

	port = <-portChan // Blockiert den Thread bis portChan in readInput gesetzt wird

	fmt.Print("TCP-Server launched...\n\n")

	// Wartet auf TCP-Verbindungen durch den port, die Bilder auf den Server hochladen
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Fatal error:\n", printTS(), err.Error())
		os.Exit(1)
	}

	// Ermöglicht (beliebig) viele TCP-Verbindungen gleichzeitig
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(printTS(), err)
			continue
		}

		// startet neue Goroutine der den verbundenen Benutzer bearbeitet
		go handleClient(conn)
	}
}

// gibt die zeit als string zurück, für error-nachrichten
func printTS() string {
	return time.Now().Format("2006-01-02_15-04-05-12345")
}
