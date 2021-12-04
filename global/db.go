package global

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var DB mongo.Database

func connectToDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal("Error connect to Db", err.Error())
	}

	DB = *client.Database(dbname)
}

func NewDBConnection(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d*performance/100)
}

func ConnectToTestDB()  {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal("Error connect to Db", err.Error())
	}

	DB = *client.Database(dbname + "_test")
}
