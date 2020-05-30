package db

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mongoPicture represents a mongo db entry
type mongoPicture struct {
	id     primitive.ObjectID `bson:"_id,omitempty"`
	Token  uuid.UUID          `bson:"token,omitempty"`
	PicID  uuid.UUID          `bson:"picId,omitempty"`
	Clicks int                `bson:"clicks"`
}
