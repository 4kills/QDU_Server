package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
)

//---------------------------------------------------------
// Web-Server-Code startet hier
//---------------------------------------------------------

var tmpl *template.Template

// Hauptfunktion des Webservers
func webServer() {
	// Wartet bis nötige variablen vom Benutzer gesetzt sind
	log.Print("Web-Server launched...\n\n")

	// assign assets to handler
	http.HandleFunc(config.DirectoryWeb, handleRequest)
	http.Handle("/pics/", http.StripPrefix("/pics/", http.FileServer(http.Dir("./pics"))))
	tmpl = template.Must(template.ParseFiles("gallery.html"))
	log.Print("Successfully assigned assets to web-server")

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
	tokstr, okMe := keys["me"]
	if ((okI || okMe) == false) || (okI && okMe) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if okMe {
		sendGallery(w, tokstr[0])
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

type user struct {
	Pics []pic
}
type pic struct {
	Name string
}

func sendGallery(w http.ResponseWriter, tokstr string) {
	tok, err := uuid.Parse(tokstr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("tokstr decode error:", err)
		return
	}

	rows, err := db.Query("SELECT pic_id FROM pics WHERE token = ?", tok[:])
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Println("db error:", err)
		return
	}

	var u user

	for rows.Next() {
		var nam16 []byte
		if err := rows.Scan(&nam16); err != nil {
			log.Println("scan row error:", err)
			continue
		}

		picID, err := uuid.FromBytes(nam16)
		if err != nil {
			log.Println("nam16 encode error:", err)
			continue
		}

		u.Pics = append(u.Pics, pic{Name: picID.String()})
	}

	tmpl.Execute(w, u)
}
