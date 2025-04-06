package scraper

import (
	"context"
    "fmt"
	// "net/http"

    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Query Methods

/*
 * Insert a course into the admin database in collection named semester
 * Arguments:
 *   courseData : Course struct with included course data
 *   semester : string representation of the semester (e.g.: Sp2025, F2025, Sp1456, etc ...)
 */
func insertCourse(courseData Course, semester string) {
	defer fmt.Printf("Inserted %s %d into %s\n", courseData.CourseCategory, courseData.CourseId, semester)
	uri := "mongodb://mongodb:27017"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	db := client.Database("admin").Collection(semester)

	if err != nil {
		fmt.Println("Error connecting to Client!")
		panic(err)
	}

	_, err = db.InsertOne(
		context.TODO(),
		courseData,
	)
	if err != nil {
		panic(err)
	}
}

// /*
//  * Establish endpoints for read operations to the database
//  */
// func DatabaseIntializer() {
// 	http.HandleFunc("/sortbyday", sortbyday);
// 	http.ListenAndServe(":8080", nil)
// }

/*
 * An exmaple insertion of how to insert a course into the db
 */
func ExampleInsertion() {
	// Example insertion:
	day := "MW"
	startTime := "0100"
	endTime := "0215"
	endTimeAMPM := "PM"
	location := "SLC"
	roomNum := 409
	info := ""

	day1 := "T"
	startTime1 := "0300"
	endTime1 := "0450"
	endTimeAMPM1 := "PM"
	location1 := "SLC"
	roomNum1 := 409
	info1 := ""

	concurrentProgrammingChild := Course{
		DeliveryMode:    "F2F",
		CourseCategory:  "CS",
		CourseId: 234,
		Section: "A",
		Crn: 10383,
		Title: "Concurrent Programming",
		Credits: 4,
		Day: &day1,
		StartTime: &startTime1,
		EndTime: &endTime1,
		EndTimeAMPM: &endTimeAMPM1,
		Location: &location1,
		RoomNum: &roomNum1,
		Instructor: "Kapolka M",
		Status: "Open",
		Limit: 20,
		Students: 14,
		Waiting: 0,
		Info: &info1,
		IsOnline: false,
	}

	concurrentProgramming := Course{
		DeliveryMode:    "F2F",
		CourseCategory:  "CS",
		CourseId: 234,
		Section: "A",
		Crn: 10383,
		Title: "Concurrent Programming",
		Credits: 4,
		Day: &day,
		StartTime: &startTime,
		EndTime: &endTime,
		EndTimeAMPM: &endTimeAMPM,
		Location: &location,
		RoomNum: &roomNum,
		Instructor: "Kapolka M",
		Status: "Open",
		Limit: 20,
		Students: 14,
		Waiting: 0,
		Info: &info,
		IsOnline: false,
		CourseChild: &concurrentProgrammingChild,
	}

	// inserts concurrentPRogramming, and concurrentProgrammingChild as a child document
	insertCourse(concurrentProgramming, "Sp2025")
}
