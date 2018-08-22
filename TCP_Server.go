package main

import (
	"bytes"
	"image/png"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/itchyny/base58-go"

	"github.com/google/uuid"
)

//---------------------------------------------------------
//TCP-Code startet hier
//---------------------------------------------------------

//Die Funktion wird für jeden verbunden Benutzer ausgeführt
func handleClient(conn net.Conn) {
	service := getService(conn)

	if service == 1 {
		sendToken(conn)
		recApproval(conn)
	}
	if service == 0 {
		sendApproval(conn)
	}

	tok58 := recToken(conn)
	sendApproval(conn)

	// erhält größe des Bildes (in Byte)
	size := recSize(conn)
	sendApproval(conn)
	// erhält das Bild, abhängig von der Größe
	buffer := recImage(conn, size)
	// erschafft einzigartigen Namen für jedes Bild
	name := createName()
	// erschafft Dateipfad für das Bild
	path := createPicPath(name)
	// erstellt die Bilddatei (hier: .png) im angegebenen Pfad
	//TODO; ADD DATABASE STUFF SOMEWHERE HERE
	draw(buffer, path)
	// schickt die URL für das Bild zurück an den Benutzer
	sendURL(conn, createURL(name))

	conn.Close()

	addToDB(name, tok58)
}

func addToDB(imgName string, tok58 []byte) {
	encoder := base58.BitcoinEncoding
	imgID, err := encoder.Decode([]byte(imgName))
	if err != nil || len(imgID) != 16 {
		log.Println("base58 encode error: imgName:", err)
	}
	tok16, err := encoder.Decode(tok58)
	if err != nil || len(tok16) != 16 {
		log.Println("base58 encode error: tok58:", err)
	}

	_, err = db.Exec("INSERT INTO pics (pic_id, token) VALUES (?, ?)", imgID, tok16)
	if err != nil {
		log.Println("db insert error:", err)
	}
}

func recToken(conn net.Conn) []byte {
	tok58 := make([]byte, 32)
	r, err := conn.Read(tok58)
	if err != nil {
		log.Println("rec token error:", err)
	}
	if r < 32 {
		log.Println("didnt fully receive token:", r)
	}
	return tok58
}

func genToken58() []byte {
	encoder := base58.BitcoinEncoding
	u := uuid.New()
	tok58, err := encoder.Encode(u[:])
	if err != nil {
		log.Println("base58 encode error:", err)
	}
	return tok58
}

func sendToken(conn net.Conn) {

	buf := genToken58()
	send, err := conn.Write(buf)
	if err != nil {
		log.Println(err)
	}
	if send != len(buf) {
		log.Println(": couldn't send token, network write error")
	}
}

func getService(conn net.Conn) byte {
	bb := make([]byte, 1)
	r, err := conn.Read(bb)
	if err != nil {
		log.Println(err)
	}
	if r < 1 {
		log.Print("didnt receive requested service properly")
	}
	return bb[0]
}

// erhält Meta-Daten wie Größe des Bildes
func recSize(conn net.Conn) int {
	// Erstellt neuen Meta-Daten buffer und liest diese vom Netzwerkstream
	bb := make([]byte, 16)
	r, err := conn.Read(bb)
	if err != nil {
		log.Println(err)
	}
	if r < 16 {
		log.Println(": didnt fully receive size: received:", r, "/16 bytes")
	}

	// Konvertiert als string verschlüsselte Größe d. Bild in einen integer
	ib := bytes.IndexByte(bb, 0)
	if ib == -1 {
		log.Println(": index overflowed: ib =", ib)
		ib = len(bb)
	}
	integer, err := strconv.Atoi(string(bb[:ib]))
	if err != nil {
		log.Println(err)
	}
	return integer
}

// Teil des Übertragungsprotokolls -> blockiert bis Client bereit
func recApproval(conn net.Conn) {
	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Println(err)
	}
	if buffer[0] != 1 {
		log.Println(": approval wasn't given")
	}
}

// Teil des Übertragungsprotokolls -> signalisiert Client, dass Server bereit
func sendApproval(conn net.Conn) {
	send, err := conn.Write([]byte{1})
	if err != nil {
		log.Println(err)
	}
	if send != 1 {
		log.Println(": couldn't send approval, network write error")
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
			log.Println(err)
		}
		rec += cur
	}
	return bytes
}

// erstellt einzigartigen Namen für jedes Bild anhand des Timestamps:
// durch die 10tel-Millisekunden ist ein doppelter Name beinahe unmöglich
func createName() string {
	return string(genToken58())
}

// schafft Dateipfad anhand des Bildnamens und standard Dateipfad
func createPicPath(name string) string {
	return filepath.Join(config.DirectoryPics, name+".png")
}

// erstellt die Bilddatei und speichert diese auf die Festplatte
func draw(buffer []byte, path string) {
	// Decoded den Bildbuffer in ein go-Image
	img, err := png.Decode(bytes.NewReader(buffer))
	if err != nil {
		log.Println(err)
		return
	}

	// Erstellt neue Datei am angegebenen Pfad und gibt ein Datei-Stream zurück
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		log.Println(err)
		return
	}

	// Schreibt das go-Image in den Datei-Stream
	errPng := png.Encode(out, img)
	if errPng != nil {
		log.Println(err)
	}
}

// Erschafft die URL anhand des Bildnamens
func createURL(name string) string {
	// TODO: improve string concatination
	dir := config.DirectoryWeb
	if config.DirectoryWeb == "" {
		dir = "/"
	}
	return "https://" + config.Domain + dir + "?i=" + name
}

// Schreibt die URL in den Netzwerk-Stream
func sendURL(conn net.Conn, s string) {
	buffer := []byte(s)
	// Schreibt Größe der URL
	if _, err := conn.Write([]byte{byte(len(buffer))}); err != nil {
		log.Println(err)
	}
	// Wartet auf Client
	recApproval(conn)
	// Schreibt URL
	if _, err := conn.Write(buffer); err != nil {
		log.Println(err)
	}
}
