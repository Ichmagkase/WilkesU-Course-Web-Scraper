package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"sync"
	"strconv"

	"wilkesu-scrapy/db"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rs/cors"
)

var mongoClient *mongo.Client

type Course struct {
	DeliveryMode   string    `json:"delivery_mode,omitempty"` // F2F; HYB; null etc.
	CourseCategory string    `json:"course_category,omitempty"` // CS; MTH; ENG etc.
	CourseId       int       `json:"course_id,omitempty"`
	Section        string    `json:"course_section,omitempty"` // INA; A; B etc.
	Crn            int       `json:"crn,omitempty"`
	Title          string    `json:"title,omitempty"`// Planet Earth; Composition; Calculus I etc.
	Credits        float32   `json:"credits,omitempty"` // 3.00; 4.00 etc.
	Day            *string   `json:"day,omitempty"`// MWF; TR; null etc.
	StartTime      *string   `json:"start_time,omitempty"`// 0100; 0800; 0430; null etc.
	EndTime        *string   `json:"end_time,omitempty"` // 0100; 0800; 0430; null etc.
	EndTimeAMPM    *string   `json:"end_time_ampm,omitempty"` // AM; PM; null.
	Location       *string   `json:"location,omitempty"`// SLC; BREIS; null etc.
	RoomNum        *int      `json:"room_num,omitempty"` // 108, 409, any number, null etc.
	Instructor     string    `json:"instructor,omitempty"` // Nye B; Simpson H; Kapolka M etc.
	Status         string    `json:"status,omitempty"` // Open; Nearly; Closed.
	Limit          int       `json:"limit,omitempty"` // Limit to number of students
	Students       int       `json:"limit,omitempty"`
	Waiting        int       `json:"waiting,omitempty"`
	Info           *string   `json:"info,omitempty"` // HONORS STUDENTS ONLY; CROSS-LISTED WITH IM 350 A etc.
	IsOnline       bool      `json:"is_online,omitempty"` // This is for full online classes (OL) not SOL or HYB
	IsCourseChild  bool      `json:"is_course_child,omitempty"` // If this course refers to a pervious course
	CourseChild    *Course
}

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
	db := mongoClient.Database("Courses").Collection(semester)

	response, err := db.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []Course
	if err = response.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	b, _ := json.Marshal(results)
	fmt.Fprintf(w, string(b))
}

/*
 * Insert a course into the admin database in collection named semester
 * Arguments:
 *   courseData : JSON byte array of course data
 *   semester : string representation of the semester (e.g.: Sp2025, F2025, Sp1456, etc ...)
 */
func InsertCourse(courseData []byte, semester string) {
	db := mongoClient.Database("Courses").Collection(semester)

	var course Course
	json.Unmarshal(courseData, &course)
	_, err := db.InsertOne(
		context.TODO(),
		course,
	)

	if err != nil {
		panic(err)

	}
}

func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there!\n")
}

func Serve() {
	var wg sync.WaitGroup
	wg.Add(1)

	// Conxonect to MongoDB instance
	mongoClient = db.Connect()

	log.Println("Initializing endpoints ...")
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

		// Serve on :8080
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
		log.Println("Endpoint being served on localhost:8080")
	}()
	wg.Wait()
}
