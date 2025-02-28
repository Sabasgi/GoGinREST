package database

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoInstance = make(map[string]*mongo.Client)
var sessionError = make(map[string]error)
var onceMap = map[string]*sync.Once{}

// GetMongoConnection method
func GetMongoConnection(mongoDsnName string) (*mongo.Client, error) {
	once, ok := onceMap[mongoDsnName]
	if !ok {
		once = &sync.Once{}
	}
	once.Do(func() {
		onceMap[mongoDsnName] = once
		clientOptions := options.Client().ApplyURI(mongoDsnName)
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Println("error connecting to mongoDB client reason:", err)
			log.Fatal(err)
		}
		// if err := client.Connect(context.TODO()); err != nil {
		// 	log.Println("error connecting to mongoDB client reason:", err)
		// 	log.Fatal(err)
		// }
		if err := client.Ping(context.TODO(), nil); err != nil {
			log.Println("error while ping mongo client reson:", err)
			log.Fatal(err)
		}
		log.Println("connected to mongoDB!!!")
		MongoInstance[mongoDsnName] = client

	})
	return MongoInstance[mongoDsnName], sessionError[mongoDsnName]
}
