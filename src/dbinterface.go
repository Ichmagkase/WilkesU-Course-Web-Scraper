package main

import (
	"context"
    "fmt"

    "go.mongodb.org/mongo-driver/v2/mongo"
    // "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

// // dbinterface main
// func dbinterface() {
// 	// uri := "mongodb://mongodb:27017"
// 	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)
// 	// opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
// 	// client, err := mongo.Connect(opts)
// 	// db := client.Database("db")

// 	// if err != nil {
// 	// 	fmt.Println("Something went wrong in the client")
// 	// 	panic(err)
// 	// }
// 	// defer func() {
// 	// 	if err = client.Disconnect(context.TODO()); err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// }()


// 	// Send a ping to confirm a successful connection
// 	var result bson.M
// 	command := bson.D{{"hello", 1}}
// 	err = db.RunCommand(context.TODO(), command).Decode(&result)
// 	if err != nil {
// 		fmt.Println("Error running command!")
// 		fmt.Println(err)
// 		panic(err)
// 		fmt.Println("Successfully panicked")
// 	}
// 	if err = client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
// 		fmt.Println("Something went wrong with pinging")
// 		panic(err)
// 	}
// 	fmt.Println("Connection established with database!")
// 	test(&db)
// }

func insertCourse(courseData course) {

}

func test() {
	uri := "mongodb://mongodb:27017"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	db := client.Database("db").Collection("courses")

	if err != nil {
		fmt.Println("Something went wrong in the client")
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	day := "MW"
	startTime := "0100"
	endTime := "0215"
	endTimeAMPM := "PM"
	location := "SLC"
	roomNum := 409
	info := ""

	day1 := "T"
	startTime1 := "0300"
	endtime1 := "0450"
	endTimeAMPM1 := "PM"
	location1 := "SLC 409"
	info1 := ""

	concurrentProgrammingChild := courseChild{
		credits: 4.00,
		day: &day1,
		startTime: &startTime1,
		endTime: &endtime1,
		endTimeAMPM: &endTimeAMPM1,
		location: &location1,
		info: &info1,
	}

	concurrentProgramming := course{
		deliveryMode:    "F2F",
		courseCategory:  "CS",
		courseId: 234,
		section: "A",
		crn: 10383,
		title: "Concurrent Programming",
		credits: 4,
		day: &day,
		startTime: &startTime,
		endTime: &endTime,
		endTimeAMPM: &endTimeAMPM,
		location: &location,
		roomNum: &roomNum,
		instructor: "Kapolka M",
		status: "Open",
		limit: 20,
		students: 14,
		waiting: 0,
		info: &info,
		isOnline: false,
		child: &concurrentProgrammingChild,
	}

	result, err := db.InsertOne(
		context.TODO(),
		concurrentProgramming,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("%s\n", result)
	fmt.Println("an insertion has been made into the table")
}
