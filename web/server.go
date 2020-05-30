package web

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"

	"github.com/4kills/base64encoding"
	"github.com/4kills/qdu_server/db"
)

var enc base64encoding.Encoder64
var tmpl *template.Template
var config configuration
var database db.Database

type configuration struct {
	directoryWeb  string
	domain        string
	portWeb       string
	directoryPics string
}

// Server starts and maintains the web server to server pics and the gallery
func Server(db db.Database) {
	database = db
	setConfig()
	// shortens links
	enc = base64encoding.New()
	// assign assets to handler
	http.HandleFunc(config.directoryWeb, handleRequest)
	http.Handle("/pics/", http.StripPrefix("/pics/", http.FileServer(http.Dir(config.directoryPics))))
	tmpl = template.Must(template.ParseFiles("gallery.html"))
	log.Print("Successfully assigned assets to web-server")

	log.Fatal("Web-Server crashed: \n\n", http.ListenAndServe(config.portWeb, nil))
}

func setConfig() {
	config = configuration{os.Getenv("WEB_PATH"), os.Getenv("DOMAIN"), os.Getenv("PORT_WEB"),
		os.Getenv("PIC_DIR")}
}

// called upon http requests
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// GET the requested picture
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

	showPic(w, pic[0])
}

func showPic(w http.ResponseWriter, picName string) {
	// writes picture into ram
	dat, err := ioutil.ReadFile(filepath.Join(config.directoryPics, picName+".png"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	// sends pic as byte stream to browser
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(dat)))
	if _, err := w.Write(dat); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	pic, err := uuidFromString(picName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("tokstrlen22 decode error: ", err)
		return
	}

	if err := database.UpdateClicks(pic, 1); err != nil {
		log.Println("db update(increment clicks) error:", err)
	}
}

type user struct {
	Pics []pic
}
type pic struct {
	Name   string
	Time   string
	Clicks int
}

func uuidFromString(tokstr string) (uuid.UUID, error) {
	var tok uuid.UUID

	b, err := enc.Decode(tokstr)
	if err != nil {
		return tok, err
	}

	tok, err = uuid.FromBytes(b[:])
	if err != nil {
		return tok, err
	}

	return tok, nil
}

func sendGallery(w http.ResponseWriter, tokstr string) {
	tok, err := uuidFromString(tokstr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("tokstrlen22 decode error: ", err)
		return
	}

	pics, err := database.QueryPics(tok)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Println("db select query error: ", err)
		return
	}

	var u user
	for _, p := range pics {
		picID := p.PictureID()
		time := p.Timestamp().UTC().Format("02-01-2006 15:04:05")
		img := pic{enc.Encode(picID[:]), time, p.PictureClicks()}
		u.Pics = append(u.Pics, img)
	}

	tmpl.Execute(w, u)
}
