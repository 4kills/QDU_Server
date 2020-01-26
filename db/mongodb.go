package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var collection *mongo.Collection

// Picture represents a mongodb entry
type Picture struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Token  uuid.UUID          `bson:"token,omitempty"`
	PicID  uuid.UUID          `bson:"picId,omitempty"`
	Clicks int                `bson:"clicks"`
}

// AddImgToDB adds the image with the specified ids to the mongo db
func AddImgToDB(imgID, tok uuid.UUID) error {
	var input Picture
	input = Picture{PicID: imgID, Token: tok}
	//input := Picture{PicID: imgID, Token: tok}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, input)
	return err
}

// QueryPics returns all the pics associated with the given user token
func QueryPics(tok uuid.UUID) ([]Picture, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{{"token", tok}})
	if err != nil {
		return []Picture{}, err
	}
	defer cur.Close(ctx)

	var pics []Picture

	for cur.Next(ctx) {
		var pic Picture
		err := cur.Decode(&pic)
		if err != nil {
			return []Picture{}, err
		}
		pics = append(pics, pic)
	}
	err = cur.Err()
	return pics, err
}

// UpdateClicks increments the click of the pic by the specified amount
func UpdateClicks(imgID uuid.UUID, amount int) error {
	filter := bson.D{{"picId", imgID}}
	update := bson.D{{"$inc", bson.D{{"clicks", amount}}}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, filter, update)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

type dbConfig struct {
	dbIP    string
	dbPort  string
	dbName  string
	colName string
}

func initDB(conf dbConfig) error {
	mongouri := fmt.Sprintf("mongodb://%s:%s", conf.dbIP, conf.dbPort)
	dbName := conf.dbName
	colName := conf.colName

	client, err := mongo.NewClient(options.Client().ApplyURI(mongouri))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("DB connection error: Ping unsuccessful: %s", err)
	}

	collection = client.Database(dbName).Collection(colName)
	return nil
}

// InitDB initializes the mongodb
func InitDB() error {
	return initDB(dbConfig{os.Getenv("dbIP"), os.Getenv("dbPort"), os.Getenv("dbName"), os.Getenv("colName")})
}
