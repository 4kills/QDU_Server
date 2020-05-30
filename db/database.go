package db

import (
	"time"

	"github.com/google/uuid"
)

type Picture interface {
	Timestamp() time.Time
	Token() uuid.UUID
	PicID() uuid.UUID
	Clicks() int
}

type Database interface {
	AddImgToDB(imgID, tok uuid.UUID) error
	QueryPics(tok uuid.UUID) ([]Picture, error)
	UpdateClicks(imgID uuid.UUID, amount int) error
	Init() error
}
