package main

import (
	"bytes"
	"image/png"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/4kills/qdu_server/db"
	"github.com/google/uuid"
)

// called upon user connect
func handleClient(conn net.Conn) {
	service := getService(conn)

	if service == 1 {
		sendToken(conn)
		recApproval(conn)
	}
	if service == 0 {
		sendApproval(conn)
	}

	tok, _ := recToken(conn)
	sendApproval(conn)

	size := recSize(conn)
	sendApproval(conn)

	buffer := recImage(conn, size)

	imgID, name := createPicID()

	path := createPicPath(name)

	draw(buffer, path)

	sendURL(conn, createURL(name))

	conn.Close()

	addToDB(imgID, tok)
}

func addToDB(imgID, tok uuid.UUID) {
	if err := db.AddImgToDB(imgID, tok); err != nil {
		log.Println("db insertion error:", err)
	}
}

func recToken(conn net.Conn) (uuid.UUID, error) {
	id := make([]byte, 36)
	_, err := conn.Read(id)
	if err != nil {
		log.Println("rec token error:", err)
	}
	if bytes.Contains(id, []byte{0}) {
		dec, err := enc.Decode(string(id[:22]))
		if err != nil {
			log.Println("rec token error:", err)
		}
		return uuid.FromBytes(dec)
	}
	return uuid.Parse(string(id))
}

func genToken() uuid.UUID {
	return uuid.New()
}

func sendToken(conn net.Conn) {

	id := genToken()
	base64 := enc.Encode(id[:])
	b := append([]byte(base64), []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}...)
	send, err := conn.Write(b)
	if err != nil {
		log.Println(err)
	}
	if send != len(id.String()) {
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

func recSize(conn net.Conn) int {
	// new metadata and read from network stream
	bb := make([]byte, 16)
	r, err := conn.Read(bb)
	if err != nil {
		log.Println(err)
	}
	if r < 16 {
		log.Println(": didnt fully receive size: received:", r, "/16 bytes")
	}

	// converts string encoded size into pic size
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

func sendApproval(conn net.Conn) {
	send, err := conn.Write([]byte{1})
	if err != nil {
		log.Println(err)
	}
	if send != 1 {
		log.Println(": couldn't send approval, network write error")
	}
}

func recImage(conn net.Conn, size int) []byte {
	bytes := make([]byte, size)

	// reads till pic fully received
	for rec := 0; rec < size; {
		cur, err := conn.Read(bytes[rec:])
		if err != nil {
			log.Println(err)
		}
		rec += cur
	}
	return bytes
}

func createPicID() (uuid.UUID, string) {
	tok := genToken()
	return tok, enc.Encode(tok[:])
}

func createPicPath(name string) string {
	return filepath.Join(os.Getenv("PIC_DIR"), name+".png")
}

func draw(buffer []byte, path string) {
	// Decodes the buffer into a go-Image
	img, err := png.Decode(bytes.NewReader(buffer))
	if err != nil {
		log.Println(err)
		return
	}

	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		log.Println(err)
		return
	}

	// writes go-img to file stream
	errPng := png.Encode(out, img)
	if errPng != nil {
		log.Println(err)
	}
}

func createURL(name string) string {
	dir := os.Getenv("WEB_PATH")
	if os.Getenv("WEB_PATH") == "" {
		dir = "/"
	}
	return "https://" + os.Getenv("DOMAIN") + dir + "?i=" + name
}

func sendURL(conn net.Conn, s string) {
	buffer := []byte(s)
	// Schreibt Größe der URL
	if _, err := conn.Write([]byte{byte(len(buffer))}); err != nil {
		log.Println(err)
	}

	recApproval(conn)

	// write URL
	if _, err := conn.Write(buffer); err != nil {
		log.Println(err)
	}
}
