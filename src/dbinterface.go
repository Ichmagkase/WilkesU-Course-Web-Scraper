package main

import (
	"context"
    "fmt"

    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

func dbinterface() {
	fmt.Println("dbinterface executed")
	uri := "mongodb://mongodb:27017"

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	db := client.Database("db")

	if err != nil {
		fmt.Println("Something went wrong in the client")
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	var result bson.M
	command := bson.D{{"hello", 1}}
	err = db.RunCommand(context.TODO(), command).Decode(&result)
	if err != nil {
		fmt.Println("Error running command!")
		fmt.Println(err)
		panic(err)
		fmt.Println("Successfully panicked")
	}
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		fmt.Println("Something went wrong with pinging")
		panic(err)
	}
	fmt.Println("Ping!")
}
