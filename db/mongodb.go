package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var collection *mongo.Collection

// AddImgToDB adds the image with the specified ids to the mongo db
func AddImgToDB(imgID, tok uuid.UUID) error {
	var input Picture
	input = Picture{PicID: imgID, Token: tok}
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
	dbIP       string
	dbPort     string
	dbName     string
	colName    string
	dbUsername string
	dbPassword string
}

func initDB(conf dbConfig) error {
	mongouri := fmt.Sprintf("mongodb://%s:%s@%s%s", conf.dbUsername, conf.dbPassword,
		conf.dbIP, conf.dbPort)
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
	return initDB(dbConfig{os.Getenv("DB_IP"), os.Getenv("PORT_DB"), os.Getenv("DB_NAME"), os.Getenv("COLL_NAME"),
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")})
}
