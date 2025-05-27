/*
 * file: start.go
 * Last updated: 5/26/2025
 * Description:
 *	 Connect to the Mongodb instance
 */
package db

import (
	"fmt"
	"log"
	"context"

	"wilkesu-scrapy/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

/*
 * Connect to MongoDB instance and return pointer to client
 */
func Connect() (*mongo.Client) {
	log.Println("Establishing MongoDB connection ...")
	var mongoClient *mongo.Client
	config := config.LoadConfig()
	uri := config.MongoUri // "mongodb://mongodb:27017"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1) // "1" is currently the only API version available
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal("config.go: ", err)
	}
	log.Println("MongoDB client initialized")
	return mongoClient
}

// For testing purposes only
func main() {
	mongoClient := Connect()

	// Ping the Database for verification
	var result bson.M
	if err := mongoClient.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		fmt.Println("Something went wrong with pinging")
		panic(err)
	}
}
