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

	"log"
	"net"
	"os"
)

// Attribute; Gültig über alle Funktionen
// Asynchron setzbar von jeder Goroutine
// (c#'s threads, co-routinen, keine sub-routinen)
var config configuration

// Haupteinstiegspunkt des Programms beim Ausführen
func main() {

	// Startet thread (goroutine) die konstant den Benutzer-Input liest und auswertet
	readInput()
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
