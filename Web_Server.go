package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/crypto/acme/autocert"

	// Biliothek zum einfacheren verteilen von http-Anfragen;
	// simuliert einen "router"
	"github.com/gorilla/mux"
)

//---------------------------------------------------------
// Web-Server-Code startet hier
//---------------------------------------------------------

// Hauptfunktion des Webservers
func webServer() {
	// Wartet bis nötige variablen vom Benutzer gesetzt sind
	fmt.Print("Web-Server launched...\n\n")

	// Adds HTTPS certificate
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("qdu.haveachin.de"),
		Cache:      autocert.DirCache("certs"),
	}

	// Konfiguriert und startet Webserver
	router := mux.NewRouter()
	server := &http.Server{
		Addr: ":1338",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
		Handler: router,
	}

	router.HandleFunc(config.DirectoryWeb, handleRequest).Methods("GET")

	go http.ListenAndServe(config.PortWeb, certManager.HTTPHandler(nil))
	if err := server.ListenAndServeTLS("", ""); err != nil {
		fmt.Println("Web-Server crashed:\n", printTS(), err)
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
		fmt.Println(printTS(), err)
	}

	// Sendet das Bild als Byte-Stream zum Broswer des Benutzers
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(dat)))
	if _, err := w.Write(dat); err != nil {
		fmt.Println(printTS(), err)
	}
}
