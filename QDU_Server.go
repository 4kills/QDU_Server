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
	"bufio"
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	// Biliothek zum einfacheren verteilen von http-Anfragen;
	// simuliert einen "router"
	"github.com/gorilla/mux"
)

// Attribute; Gültig über alle Funktionen
// Asynchron setzbar von jeder Goroutine
// (c#'s threads, co-routinen, keine sub-routinen)
var domain string
var directory string
var directoryWeb string
var port string
var portWeb string

// Haupteinstigspunkt des Programms beim Ausführen
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
		fmt.Println("Fatal error:\n", err.Error())
		os.Exit(1)
	}

	// Ermöglicht (beliebig) viele TCP-Verbindungen gleichzeitig
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		// startet neue Goroutine der den verbundenen Benutzer bearbeitet
		go handleClient(conn)
	}
}

//---------------------------------------------------------
// Web-Server-Code startet hier
//---------------------------------------------------------

// Hauptfunktion des Webservers
func webServer(pC, dirC <-chan string) {
	// Wartet bis nötige variablen vom Benutzer gesetzt sind
	select {
	case portWeb = <-pC:
		directoryWeb = <-dirC
	case directoryWeb = <-dirC:
		portWeb = <-pC
	}
	fmt.Print("Web-Server launched...\n\n")

	// Konfiguriert und startet Webserver
	router := mux.NewRouter()
	router.HandleFunc(directoryWeb, handleRequest).Methods("GET")
	if err := http.ListenAndServe(portWeb, router); err != nil {
		fmt.Println("Web-Server crashed:\n", err)
		os.Exit(1)
	}
}

// Die Funktion die aufgerufen wird, wenn eine http-Anfrage hereinkommt
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Liest aus der URL durch die GET-Methode das angefragte Bild aus
	keys := r.URL.Query()
	pic := keys["i"]
	if len(pic) < 1 {
		return
	}
	// schreibt kompletten inhalt der Bild-Datei in den RAM
	dat, err := ioutil.ReadFile(filepath.Join(directory, pic[0]+".png"))
	if err != nil {
		fmt.Println(err)
	}

	// Sendet das Bild als Byte-Stream zum Broswer des Benutzers
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(dat)))
	if _, err := w.Write(dat); err != nil {
		fmt.Println(err)
	}
}

//---------------------------------------------------------
//TCP-Code startet hier
//---------------------------------------------------------

//Die Funktion wird für jeden verbunden Benutzer ausgeführt
func handleClient(conn net.Conn) {
	defer conn.Close()
	// erhält größe des Bildes (in Byte)
	size := recMetaData(conn)
	sendApproval(conn)
	// erhält das Bild, abhängig von der Größe
	buffer := recImage(conn, size)
	// erschafft einzigartigen Namen für jedes Bild
	name := createName()
	// erschafft Dateipfad für das Bild
	path := createPicPath(name)
	// erstellt die Bilddatei (hier: .png) im angegebenen Pfad
	draw(buffer, path)
	// schickt die URL für das Bild zurück an den Benutzer
	sendURL(conn, createURL(name))
}

// erhält Meta-Daten wie Größe des Bildes
func recMetaData(conn net.Conn) int {
	// Erstellt neuen Meta-Daten buffer und liest diese vom Netzwerkstream
	bytes := make([]byte, 16)
	_, err := conn.Read(bytes)
	if err != nil {
		fmt.Println(err)
	}

	// Konvertiert als string verschlüsselte Größe d. Bild in einen integer
	s := string(bytes)
	integer, err2 := strconv.Atoi(s[:strings.IndexByte(s, 0)])
	if err2 != nil {
		fmt.Println(err2)
	}
	return integer
}

// Teil des Übertragungsprotokolls -> blockiert bis Client bereit
func recApproval(conn net.Conn) {
	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}
	if buffer[0] != 1 {
		fmt.Println("approval wasn't given")
	}
}

// Teil des Übertragungsprotokolls -> signalisiert Client, dass Server bereit
func sendApproval(conn net.Conn) {
	send, err := conn.Write([]byte{1})
	if err != nil {
		fmt.Println(err)
	}
	if send != 1 {
		fmt.Println("couldn't send approval, network write error")
	}
}

// erhält Bild
func recImage(conn net.Conn, size int) []byte {
	// Erstellt go-slice (mischung aus arrays und listen in go)
	// mit größe des Bildes als byte-buffer
	bytes := make([]byte, size)
	// liest so lange bytes vom Netzwerk-Stream
	// bis das Bild vollständig angekommen ist
	for rec := 0; rec < size; {
		cur, err := conn.Read(bytes[rec:])
		if err != nil {
			fmt.Println(err)
		}
		rec += cur
	}
	return bytes
}

// erstellt einzigartigen Namen für jedes Bild anhand des Timestamps:
// durch die 10tel-Millisekunden ist ein doppelter Name beinahe unmöglich
func createName() string {
	// Erhält Timestamp als string und kürzt diesen
	str := strings.Replace(time.Now().Format("2006-01-02_15-04-05-12345"), "-", "", -1)
	str = strings.Replace(str, "_", "", -1)
	str = strings.Replace(str, "0", "", -1)
	// Konvertierung des strings in eine Zahl und base58 encoding (um URL zu kürzen)
	num, err := strconv.Atoi(str)
	if err != nil || num < 0 {
		fmt.Println(err)
		fmt.Println("if err = null -> num is negative")
	}
	return base58Encoding(int64(num))
}

// Kodiert eine zahl zu einer Basis 58-Zahl um URL zu kürzen
// wie Hexadezimal nur zur Basis 58; Basis 64 ungeeignet wegen "/"-Sonderzeichen
func base58Encoding(num int64) string {
	const codeSet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	base := int64(len(codeSet)) // 58
	var encoded []string

	// itteriert durch die Zahl und hängt das nötige base58-Zeichen an output-string
	for num != 0 {
		remainder := num % base
		encoded = append([]string{codeSet[remainder : remainder+1]}, strings.Join(encoded, ""))
		num = num / base
	}

	return strings.Join(encoded, "")
}

// schafft Dateipfad anhand des Bildnamens und standard Dateipfad
func createPicPath(name string) string {
	return filepath.Join(directory, name+".png")
}

// erstellt die Bilddatei und speichert diese auf die Festplatte
func draw(buffer []byte, path string) {
	// Decoded den Bildbuffer in ein go-Image
	img, errImg := png.Decode(bytes.NewReader(buffer))
	if errImg != nil {
		fmt.Println(errImg)
		return
	}

	// Erstellt neue Datei am angegebenen Pfad und gibt ein Datei-Stream zurück
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Schreibt das go-Image in den Datei-Stream
	errPng := png.Encode(out, img)
	if errPng != nil {
		fmt.Println(err)
	}
}

// Erschafft die URL anhand des Bildnamens
func createURL(name string) string {
	// TODO: improve string concatination
	dir := directoryWeb
	if directoryWeb == "" {
		dir = "/"
	}
	return "http://" + domain + dir + "?i=" + name
}

// Schreibt die URL in den Netzwerk-Stream
func sendURL(conn net.Conn, s string) {
	buffer := []byte(s)
	// Schreibt Größe der URL
	if _, err := conn.Write([]byte{byte(len(buffer))}); err != nil {
		fmt.Println(err)
	}
	// Wartet auf Client
	recApproval(conn)
	// Schreibt URL
	if _, err := conn.Write(buffer); err != nil {
		fmt.Println(err)
	}
}

//---------------------------------------------------------
// Lesen und Auswerten von User-Input startet hier
//---------------------------------------------------------

// Haupteinstiegsfunktion für das Lesen von User-Input
func readInput(c, tcpC, dirC chan<- string) {
	fmt.Print("QDU-Server launched...\n",
		"Launch servers by entering the respetive settings\n\n")

	fmt.Println("Write '/help' for a list of commands")
	fmt.Print("-------------------------------------\n\n")
	// liest konstant Konsolenzeile
	for {
		readLine(c, tcpC, dirC)
	}
}

// Liest Konsolenzeile und lässt diese auf Befehle überprüfen
func readLine(c, tcpC, dirC chan<- string) {
	// Liest Konsolen-Stream und bricht bei "Enter" ab zu lesen
	const inputDelimiter = '\n'
	r := bufio.NewReader(os.Stdin)
	input, err := r.ReadString(inputDelimiter)
	if err != nil {
		fmt.Println(err)
		return
	}

	// return bei unnötigem Input
	if !strings.Contains(input, "/") {
		return
	}

	// Formatiert input
	input = strings.Replace(input, "\n", "", -1)
	input = strings.Trim(input, " ")
	if runtime.GOOS == "windows" {
		input = input[:len(input)-1]
	}
	input = strings.Replace(input, " ", "", -1)

	// überprüft auf Befehle
	checkCommands(input, c, tcpC, dirC)
}

// Überprüft auf Befehle und setzt die Channel
func checkCommands(s string, c, tcpC, dirC chan<- string) {
	// TODO: reduce hard coded strings (via export?)

	// switch in Go kann benutzt werden um lange if-else
	// ausdrücke leserlicher zu gestalten.
	switch {
	case len(s) > 11 && (s[:11] == "/setPortTCP"):
		c <- ":" + s[11:]
		fmt.Print("successfully set port\n\n")
	case s == "/getPortTCP":
		if len(port) < 1 {
			fmt.Print("\n\n")
			break
		}
		fmt.Print(port[1:], "\n\n")
	case len(s) > 11 && (s[:11] == "/setPortWeb"):
		tcpC <- ":" + s[11:]
		fmt.Print("successfully set portWeb\n\n")
	case s == "/getPortWeb":
		if len(portWeb) < 1 {
			fmt.Print("\n\n")
			break
		}
		fmt.Print(portWeb[1:], "\n\n")
	case s == "/getDirWeb":
		fmt.Print(directoryWeb)
	case len(s) > 10 && (s[:10] == "/setDirWeb"):
		dirC <- s[10:]
		fmt.Print("successfully set webDir\n\n")
	case s == "/help":
		fmt.Print("\n-------------------------------------------------------------------------------",
			"\n/help \t\t\t\t- returns a list of commands",
			"\n/getDomain \t\t\t- gets current domain for links",
			"\n/setDomain <domain>\t\t- sets domain for returning link",
			"\n/getDir \t\t\t- gets output directory of pics",
			"\n/setDir <abs path> \t\t- sets output directory of pics",
			"\n/getDirWeb \t\t\t- gets directory for web requests",
			"\n/setDirWeb <rel path>\t\t- sets directory for web requests; may be empty",
			"\n/getPortTCP \t\t\t- gets the admin defined port",
			"\n/setPortTCP <port> \t\t- sets the tcp port. Set this first",
			"\n/getPortWeb \t\t\t- gets the web port",
			"\n/setPortWeb <port>\t\t- sets the web port. Set this second",
			"\n/quit \t\t\t\t- exits the application w/ exit status 3",
			"\n/info \t\t\t\t- returns information about the programm",
			"\n-------------------------------------------------------------------------------\n\n")
	case s == "/getDir":
		fmt.Println(directory)
	case len(s) > 7 && (s[:7] == "/setDir"):
		directory = s[7:]
		fmt.Print("successfully set output directory\n\n")
	case s == "/getDomain":
		fmt.Println(domain)
	case len(s) > 10 && (s[:10] == "/setDomain"):
		domain = s[10:]
		fmt.Print("successfully set domain\n\n")
	case s == "/quit":
		os.Exit(3)
	case s == "/info":
		fmt.Print("\n-------------------------------------------------------------------------------",
			"\nDevelopment Framework: Golang 1.10",
			"\nVersion: 1.0 (stable release)",
			"\nCreator: Dominik Ochs",
			"\n-------------------------------------------------------------------------------\n\n")
	}
}
