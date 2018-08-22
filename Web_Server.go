package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/itchyny/base58-go"
)

//---------------------------------------------------------
// Web-Server-Code startet hier
//---------------------------------------------------------

// Hauptfunktion des Webservers
func webServer() {
	// Wartet bis nötige variablen vom Benutzer gesetzt sind
	log.Print("Web-Server launched...\n\n")

	// assign assets to handler
	http.HandleFunc(config.DirectoryWeb, handleRequest)
	http.Handle("/pics/", http.StripPrefix("/pics/", http.FileServer(http.Dir("./pics"))))

	go http.ListenAndServe(":http", nil)
	log.Fatal("Web-Server crashed: \n\n", http.ListenAndServeTLS(config.PortWeb,
		"/etc/letsencrypt/live/haveachin.de/fullchain.pem",
		"/etc/letsencrypt/live/haveachin.de/privkey.pem", nil))
}

// Die Funktion die aufgerufen wird, wenn eine http-Anfrage hereinkommt
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Liest aus der URL durch die GET-Methode das angefragte Bild aus
	keys := r.URL.Query()
	pic, okI := keys["i"]
	tok58, okMe := keys["me"]
	if okI || okMe == false || okI && okMe == true {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if len(pic) < 16 || len(tok58) < 16 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if okMe {
		sendGallery(w, tok58)
		return
	}

	// schreibt kompletten inhalt der Bild-Datei in den RAM
	dat, err := ioutil.ReadFile(filepath.Join(config.DirectoryPics, pic[0]+".png"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	// Sendet das Bild als Byte-Stream zum Broswer des Benutzers
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(dat)))
	if _, err := w.Write(dat); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}

func sendGallery(w http.ResponseWriter, tok58 []string) {
	encoder := base58.BitcoinEncoding
	tok16, err := encoder.Decode([]byte(tok58[0]))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("tok58 decode error:", err)
		return
	}

	rows, err := db.Query("SELECT pic_id FROM pics WHERE token = ?", tok16)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Println("db error:", err)
		return
	}

	var links []string
	for rows.Next() {
		var nam16 []byte
		if err := rows.Scan(nam16); err != nil {
			log.Println("scan row error:", err)
			continue
		}
		nam58, err := encoder.Encode(nam16)
		if err != nil {
			log.Println("nam16 encode error:", err)
			continue
		}

		links = append(links, string(nam58))
	}

	w.WriteHeader(http.StatusOK)
	for _, ele := range links {
		if _, err := fmt.Fprintf(w, "<img src='pics/"+ele+".png' alt='thepic' style='width:320px;height:235px;'>"); err != nil {
			log.Println(err)
			continue
		}
	}
}
