package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

//---------------------------------------------------------
// Web-Server-Code startet hier
//---------------------------------------------------------

// Hauptfunktion des Webservers
func webServer() {
	// Wartet bis nötige variablen vom Benutzer gesetzt sind
	log.Print("Web-Server launched...\n\n")

	// Adds HTTPS certificate
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Domain),
		Cache:      autocert.DirCache("certs"),
	}

	// Konfiguriert und startet Webserver
	router := mux.NewRouter()
	server := &http.Server{
		Addr: config.PortWeb,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
		Handler: router,
	}

	router.HandleFunc(config.DirectoryWeb, handleRequest).Methods("GET")

	go http.ListenAndServe(config.PortWeb, certManager.HTTPHandler(nil))
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Println("Web-Server crashed:\n", err)
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
