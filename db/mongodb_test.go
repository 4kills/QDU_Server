package db

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
)

func TestDB(t *testing.T) {
	db := mongoDB{}
	err := db.init(dbConfig{dbIP: "192.168.178.25", dbPort: "27017", dbName: "QDU", colName: "pics"})
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	img, _ := uuid.FromBytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF})
	tok, _ := uuid.FromBytes([]byte{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1})
	err = db.AddImgToDB(img, tok)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	err = db.UpdateClicks(img, 1)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	pics, err := db.QueryPics(tok)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	for _, pic := range pics {
		fmt.Println(pic)
	}
}
