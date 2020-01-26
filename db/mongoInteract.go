package db

import (
	"context"
	"fmt"
	"log"
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
var picID = "picId"
var token = "token"
var clicks = "clicks"

// AddImgToDB adds the image with the specified ids to the mongo db
func AddImgToDB(imgID, tok uuid.UUID) error {
	input := bson.M{picID: imgID, token: tok, clicks: 0}
	return insert(input)
}

func insert(input bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, input)
	return err
}

// Picture represents a picture in a mongodb entry
type Picture struct {
	ID     primitive.ObjectID
	PicID  uuid.UUID
	Clicks int
}

// QueryPics returns all the pics associated with the given user token
func QueryPics(tok uuid.UUID) ([]Picture, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.M{token: tok})
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
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return pics, nil
}

// UpdateClicks increments the click of the pic by the specified amount
func UpdateClicks(imgID uuid.UUID, amount int) error {
	filter := bson.M{picID: imgID}
	opt := options.FindOneAndUpdate().SetUpsert(false)
	update := bson.M{"$inc": bson.M{clicks: amount}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var updatedDoc bson.M
	err := collection.FindOneAndUpdate(ctx, filter, update, opt).Decode(&updatedDoc)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// InitDB initializes the mongodb
func InitDB() {
	mongouri := fmt.Sprintf("mongodb://%s:%s", os.Getenv("dbIP"), os.Getenv("dbPort"))
	dbName := os.Getenv("dbName")
	colName := os.Getenv("colName")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongouri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("DB connection error: Ping unsuccessful: ", err)
	}

	collection = client.Database(dbName).Collection(colName)
}
