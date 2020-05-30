package db

import (
	"time"

	"github.com/google/uuid"
)

// Picture represents an Interface of the Pictures as saved in the database
type Picture interface {
	Timestamp() time.Time
	UserToken() uuid.UUID
	PictureID() uuid.UUID
	PictureClicks() int
}

// Database represents a database connection with methods used for handling Pictures
type Database interface {
	AddImgToDB(imgID, tok uuid.UUID) error
	QueryPics(tok uuid.UUID) ([]Picture, error)
	UpdateClicks(imgID uuid.UUID, amount int) error
	Init() error
}

// New returns the current implementation of a Database
func New() Database {
	return mongoDB{}
}

type dbConfig struct {
	dbIP       string
	dbPort     string
	dbName     string
	colName    string
	dbUsername string
	dbPassword string
}
