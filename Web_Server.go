package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
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

	go http.ListenAndServe(":http", nil)
	log.Fatal("Web-Server crashed: \n\n", http.ListenAndServeTLS(config.PortWeb,
		"/etc/letsencrypt/live/haveachin.de/fullchain.pem",
		"/etc/letsencrypt/live/haveachin.de/privkey.pem", nil))
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
	dat, err := ioutil.ReadFile(filepath.Join(config.DirectoryPics, pic[0]+".png"))
	if err != nil {
		log.Println(err)
	}

	// Sendet das Bild als Byte-Stream zum Broswer des Benutzers
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(dat)))
	if _, err := w.Write(dat); err != nil {
		log.Println(err)
	}
}
