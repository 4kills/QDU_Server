package main

import (
	"bytes"
	"fmt"
	"image/png"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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
		fmt.Println(printTS(), err)
	}

	// Konvertiert als string verschlüsselte Größe d. Bild in einen integer
	s := string(bytes)
	integer, err2 := strconv.Atoi(s[:strings.IndexByte(s, 0)])
	if err2 != nil {
		fmt.Println(printTS(), err2)
	}
	return integer
}

// Teil des Übertragungsprotokolls -> blockiert bis Client bereit
func recApproval(conn net.Conn) {
	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(printTS(), err)
	}
	if buffer[0] != 1 {
		fmt.Println(printTS(), ": approval wasn't given")
	}
}

// Teil des Übertragungsprotokolls -> signalisiert Client, dass Server bereit
func sendApproval(conn net.Conn) {
	send, err := conn.Write([]byte{1})
	if err != nil {
		fmt.Println(printTS(), err)
	}
	if send != 1 {
		fmt.Println(printTS(), ": couldn't send approval, network write error")
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
			fmt.Println(printTS(), err)
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
	str = str[1 : len(str)-3]

	// Konvertierung des strings in eine Zahl und base58 encoding (um URL zu kürzen)
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil || num < 0 {
		fmt.Println(printTS(), err)
		fmt.Println("if err = null -> num is negative")
	}
	return base58Encoding(num)
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
	return filepath.Join(config.DirectoryPics, name+".png")
}

// erstellt die Bilddatei und speichert diese auf die Festplatte
func draw(buffer []byte, path string) {
	// Decoded den Bildbuffer in ein go-Image
	img, errImg := png.Decode(bytes.NewReader(buffer))
	if errImg != nil {
		fmt.Println(printTS(), errImg)
		return
	}

	// Erstellt neue Datei am angegebenen Pfad und gibt ein Datei-Stream zurück
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		fmt.Println(printTS(), err)
		return
	}

	// Schreibt das go-Image in den Datei-Stream
	errPng := png.Encode(out, img)
	if errPng != nil {
		fmt.Println(printTS(), err)
	}
}

// Erschafft die URL anhand des Bildnamens
func createURL(name string) string {
	// TODO: improve string concatination
	dir := config.DirectoryWeb
	if config.DirectoryWeb == "" {
		dir = "/"
	}
	return "http://" + config.Domain + dir + "?i=" + name
}

// Schreibt die URL in den Netzwerk-Stream
func sendURL(conn net.Conn, s string) {
	buffer := []byte(s)
	// Schreibt Größe der URL
	if _, err := conn.Write([]byte{byte(len(buffer))}); err != nil {
		fmt.Println(printTS(), err)
	}
	// Wartet auf Client
	recApproval(conn)
	// Schreibt URL
	if _, err := conn.Write(buffer); err != nil {
		fmt.Println(printTS(), err)
	}
}
