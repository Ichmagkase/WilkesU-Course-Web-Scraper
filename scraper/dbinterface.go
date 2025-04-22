package scraper

import (
	"context"
    "fmt"
	"net/http"
	"sync"

    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var mongoClient *mongo.Client

func responseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("MONGO CLIENT: (responseHandler): ")
	fmt.Println(mongoClient)
	// get responses
	params := r.URL.Query()
	semester := params["semester"][0]
	// sortby := params["sort"][0]
	deliveryMode := params["dm"][0]
	// category := params["category"][0]
	// section := params["section"][0]
	// credits := params["credits"][0]
	// location := params["loc"][0]
	instructor := params["instructor"][0]
	// status := params["status"][0]

	// Set up db connection
	db := mongoClient.Database("admin").Collection(semester)

	filter := bson.D{
		{"instructor",
			bson.D{
				{"$regex", instructor},
				{"$options", "i"},
			},
		},
		{"dm",
			bson.D{
				{"$regex", deliveryMode},
				{"$options", "i"},
			},
		},
	}
	response, err := db.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []Course
	if err = response.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for i := 0; i < len(results); i++ {
		fmt.Fprintf(w, "%s\n", results[i].Title)
	}
}

func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there!\n")
}

 /*
  * Establish endpoints for read operations to the database
  */
func DatabaseIntializer() {
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("DatabaseInitializer()")
	fmt.Println("Database is being initialized!")
	go func() {
		defer wg.Done()
		// Run database
		uri := "mongodb://mongodb:27017"
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
		mongoClient, _ = mongo.Connect(opts)

		fmt.Printf("MONGO CLIENT (DatabaseInitializer): ")
		fmt.Println(mongoClient)

		// Keep the Go compiler happy
		var result bson.M
		if err := mongoClient.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
			fmt.Println("Something went wrong with pinging")
			panic(err)
		}
	}()

	fmt.Println("Endpoints are being initialized!")
	go func() {
		defer wg.Done()
		// Run HTTP server
		mux := http.NewServeMux()
		mux.HandleFunc("/filter", responseHandler)
		mux.HandleFunc("/test", testResponse)
		server := http.Server{
			Addr: ":8080",
			Handler: mux,
		}
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}

// Query Methods

/*
 * Insert a course into the admin database in collection named semester
 * Arguments:
 *   courseData : Course struct with included course data
 *   semester : string representation of the semester (e.g.: Sp2025, F2025, Sp1456, etc ...)
 */
func insertCourse(courseData Course, semester string) {
	fmt.Printf("MONGO CLIENT (insertCourse): ")
	fmt.Println(mongoClient)
	db := mongoClient.Database("admin").Collection(semester)

	_, err := db.InsertOne(
		context.TODO(),
		courseData,
	)
	if err != nil {
		panic(err)
	}
}
