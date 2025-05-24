package scraper

import (
	"context"
    "fmt"
	"net/http"
	"encoding/json"
	"sync"
	"strconv"

    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rs/cors"
)

var mongoClient *mongo.Client

func responseHandler(w http.ResponseWriter, r *http.Request) {

	// Handle CORS middleware Preflight options
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get responses
	params := r.URL.Query()
	semester := params["semester"][0]
	deliverymode := params["deliverymode"]
	category := params["category"]
	credits := params["credits"]
	location := params["location"]
	instructor := params["instructor"]
	status := params["status"]
	crn := params["crn"]

	filterOpts := []string{
		"deliverymode",
		"coursecategory",
		"location",
		"instructor",
		"status",
	}

	receivedParams := [][]string {
		deliverymode,
		category,
		location,
		instructor,
		status,
	}

	fmt.Println(receivedParams)

	// Handle string parameters
	filter := bson.D{}
	for i, _  := range receivedParams {
		if len(receivedParams[i]) > 0 {
			filter = append(filter, bson.E{
				filterOpts[i],
				bson.D{
					{"$regex", receivedParams[i][0]},
					{"$options", "i"},
				},
			})
		}
	}

	// Handle integer parameters
	if len(credits) > 0 {
		creditsint, _ := strconv.Atoi(credits[0])
		filter = append(filter, bson.E{
			"credits", creditsint,
		})
	}

	if len(crn) > 0 {
		crnint, _ := strconv.Atoi(crn[0])
		filter = append(filter, bson.E{
			"crn", crnint,
		})
	}

	// Set up db connection
	db := mongoClient.Database("admin").Collection(semester)

	response, err := db.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []Course
	if err = response.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	b, err := json.Marshal(results)
	fmt.Fprintf(w, string(b))

	// for i := 0; i < len(results); i++ {
	// 	fmt.Fprintf(w, "%s %s\n", results[i].Title, results[i].Instructor)
	// }
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

		// Include CORS headers
		c := cors.New(cors.Options{
			AllowedOrigins: []string{
				"http://localhost:5174",
			},
			AllowedMethods: []string{
				"GET",
			},
			AllowedHeaders: []string{
				"*",
			},
			AllowCredentials: true,
		})

		// Build server options
		mux := http.NewServeMux()
		mux.HandleFunc("/filter", responseHandler)
		mux.HandleFunc("/test", testResponse)
		handler := c.Handler(mux)
		server := http.Server{
			Addr: ":8080",
			Handler: handler,
		}

		// Serve on server
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
