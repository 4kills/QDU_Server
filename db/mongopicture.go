package db

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mongoPicture represents a mongo db entry
type mongoPicture struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Token  uuid.UUID          `bson:"token,omitempty"`
	PicID  uuid.UUID          `bson:"picId,omitempty"`
	Clicks int                `bson:"clicks"`
}

func (m mongoPicture) Timestamp() time.Time {
	return m.ID.Timestamp()
}

func (m mongoPicture) UserToken() uuid.UUID {
	return m.Token
}

func (m mongoPicture) PictureID() uuid.UUID {
	return m.PicID
}

func (m mongoPicture) PictureClicks() int {
	return m.Clicks
}
