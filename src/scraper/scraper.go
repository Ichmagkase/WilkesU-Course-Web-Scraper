package scraper

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"strconv"
	"golang.org/x/net/html"
	"errors"
	"os"
	"regexp"
	"sync"
	"log"
	"context"
	"sync/atomic"
	"runtime/trace"
	"wilkesu-scrapy/api"
)

/* The course struct is what a course is expected to look like.

Any variable that is a pointer to a type is an optional paramteter.
*/
type Course struct {
	DeliveryMode string `json:"delivery_mode,omitempty"` // F2F; HYB; null etc.
	CourseCategory string `json:"course_category,omitempty"` // CS; MTH; ENG etc.
	CourseId int `json:"course_id,omitempty"`
	Section string `json:"course_section,omitempty"` // INA; A; B etc.
	Crn int `json:"crn,omitempty"`
	Title string `json:"title,omitempty"`// Planet Earth; Composition; Calculus I etc.
	Credits float32 `json:"credits,omitempty"` // 3.00; 4.00 etc.
	Day *string `json:"day,omitempty"`// MWF; TR; null etc.
	StartTime *string `json:"start_time,omitempty"`// 0100; 0800; 0430; null etc.
	EndTime *string `json:"end_time,omitempty"` // 0100; 0800; 0430; null etc.
	EndTimeAMPM *string `json:"end_time_ampm,omitempty"` // AM; PM; null.
	Location *string `json:"location,omitempty"`// SLC; BREIS; null etc.
	RoomNum *int `json:"room_num,omitempty"` // 108, 409, any number, null etc.
	Instructor string `json:"instructor,omitempty"` // Nye B; Simpson H; Kapolka M etc.
	Status string `json:"status,omitempty"` // Open; Nearly; Closed.
	Limit int `json:"limit,omitempty"` // Limit to number of students
	Students int `json:"students,omitempty"`
	Waiting int `json:"waiting,omitempty"`
	Info *string `json:"info,omitempty"` // HONORS STUDENTS ONLY; CROSS-LISTED WITH IM 350 A etc.
	IsOnline bool `json:"is_online,omitempty"` // This is for full online classes (OL) not SOL or HYB
	IsCourseChild bool `json:"is_course_child,omitempty"` // If this course refers to a pervious course
	CourseChild *Course
}

/* Parsing functions */
type fieldFunc func (*Course, *html.Tokenizer, *int, html.Token) error

func courseToString(c Course) string {
	/* courseToString takes a course and parses it into a string

	Arguments:
		c (course): the course to stringify.

	Returns:
		string: string representation of a course.
	*/

    safeString := func(s *string) string {
        if s == nil {
            return "N/A"
        }
        return *s
    }

    // Helper function to safely dereference pointer integers
    safeInt := func(i *int) string {
        if i == nil {
            return "N/A"
        }
        return fmt.Sprintf("%d", *i)
    }

    // Build the course description
    description := fmt.Sprintf(`Course Details:
	- Course ID: %d
	- Title: %s
	- Delivery Mode: %s
	- Category: %s
	- Section: %s
	- CRN: %d
	- Credits: %.2f
	- Instructor: %s
	- Status: %s

	Schedule:
	- Days: %s
	- Start Time: %s
	- End Time: %s %s

	Enrollment:
	- Limit: %d
	- Current Students: %d
	- Waiting List: %d

	Additional Information:
	- Location: %s
	- Room Number: %s
	- Online: %t
	- Special Info: %s`,
        c.CourseId,
        c.Title,
        c.DeliveryMode,
        c.CourseCategory,
        c.Section,
        c.Crn,
        c.Credits,
        c.Instructor,
        c.Status,
        safeString(c.Day),
        safeString(c.StartTime),
        safeString(c.EndTime),
        safeString(c.EndTimeAMPM),
        c.Limit,
        c.Students,
        c.Waiting,
        safeString(c.Location),
        safeInt(c.RoomNum),
        c.IsOnline,
        safeString(c.Info))

    if c.CourseChild != nil {
        childString := fmt.Sprintf(`
	Additional Course Section:
	%s`, courseToString(*c.CourseChild))
        
        description += childString
    }

    return description
}

func getDeliveryMode(c *Course, tokenizer *html.Tokenizer, fieldCount *int, startToken html.Token) error {
	/* getDeliveryMode gets the delivery mode from the course.

	Delivery modes is how the course is offered (F2F, HYB, OL etc.)

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Delivery Mode field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Delivery Mode Token found: %s\n", token.Data)
			c.DeliveryMode = token.Data
		}
	}
	return nil
}

func getCourseCategoryAndId(c *Course, tokenizer *html.Tokenizer, fieldCount *int, startToken html.Token) error {
	/* getCourseCategoryAndId gets the course category 
	and id of the course.

	Course category is where the course belongs to (CS, MTH, ENG etc.)
	Course id is the number of that course. 

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Course Category field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Course Category and Id found: %s\n", token.Data)
			splitData := strings.Split(token.Data, " ");
			if (len(splitData) != 2) {
				return errors.New(fmt.Sprintf("Course category and id in unexpected format." + 
									   " Got %d; Expected 2.", len(splitData)))
			}
			c.CourseCategory = splitData[0]
			courseId, err := strconv.Atoi(splitData[1])
			if err != nil {
				return err
			}
			c.CourseId = courseId
		}
	}
	return nil
}

func getSection(c *Course, tokenizer *html.Tokenizer, fieldCount *int, startToken html.Token) error {
	/* getSection gets the section from the course.

	Section is the group of the Section (A, B, INA etc.)

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Section field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Section Token found: %s\n", token.Data)
			c.Section = token.Data
		}
	}
	return nil
}

func getCRN(c *Course, tokenizer *html.Tokenizer, fieldCount *int,  startToken html.Token) error {
	/* getCRN gets the CRN from the course.

	CRN is the course registration number (31233,23213, 0 etc.)

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of CRN field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("CRN Token found: %s\n", token.Data)
			parsedCRN, err := strconv.Atoi(token.Data)
			if err != nil {
				return err
			}
			c.Crn = parsedCRN
		}
	}
	return nil
}

func getTitle(c *Course, tokenizer *html.Tokenizer, fieldCount *int,  startToken html.Token) error {
	/* getTitle gets the name from the course.

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Title field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Title Token found: %s\n", token.Data)
			c.Title = token.Data
		}
	}
	return nil
}

func getCredits(c *Course, tokenizer *html.Tokenizer, fieldCount *int, startToken html.Token) error {
	/* getCredits gets the credits from the course.

	Credits are floats (3.00, 4.00, 1.00, etc.)

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Credits field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Credits Token found: %s\n", token.Data)
			parsedCredits, err := strconv.ParseFloat(token.Data, 32)
			if err != nil {
				return err
			}
			c.Credits = float32(parsedCredits)
		}
	}
	return nil
}

func getDay(c *Course, tokenizer *html.Tokenizer, fieldCount *int,  startToken html.Token) error {
	/* getDay gets the days given from the course.

	Days are formated as TR, MWF, WF, R etc.

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	// Check if date, time, and location exists
	for _, attr := range startToken.Attr {
		// Day, Time, Location, DNE
		if (attr.Key == "colspan" && attr.Val == "3") {

			// Move tokenzier to the next row
			fmt.Println("Day, Time, and Location unknow . . . skipping to instructor")
			tokenizer.Next()
			t := tokenizer.Token()
			c.Location = &t.Data
			tokenizer.Next()
			*fieldCount = *fieldCount + 2
			return nil
		}
	}

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Day field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Day Token found: %s\n", token.Data)
			c.Day = &token.Data
		}
	}
	return nil
}

func getTime(c *Course, tokenizer *html.Tokenizer, fieldCount *int,  startToken html.Token) error {
	/* getTime gets the time given from the course.

	Course time is formatted as such:
	0900-0950AM; 1000-1150AM; 1100-1250PM etc.
	
	Course time is divided into serveral fields:
	startTime *string // 0100; 0800; 0430; null etc.
	EndTime *string // 0100; 0800; 0430; null etc.
	EndTimeAMPM *string // AM; PM; null.

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Time field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Time Token found: %s\n", token.Data)

			timeFormat := "^([0-9]{4}-[0-9]{4})(?:AM|PM)$"
			re := regexp.MustCompile(timeFormat)
			if (re.Match([]byte(token.Data))) {
				time := token.Data

				// Start time: XX:XX 
				formattedStartTime := fmt.Sprintf("%s:%s",time[0:2],time[2:4]) 
				c.StartTime = &formattedStartTime

				// End time: XX:XX 
				formattedEndTime := fmt.Sprintf("%s:%s",time[5:7], time[7:9])
				c.EndTime = &formattedEndTime

				// End time AMPM: AM or PM
				amOrPm := time[9:11]
				c.EndTimeAMPM = &amOrPm

			} else {
				if (token.Data == "TBA") {
					c.StartTime = &token.Data
					c.EndTime = &token.Data
					c.EndTimeAMPM = &token.Data
				} else {
					return errors.New(fmt.Sprintf("Error: time was not in the right format. Got %s, Expected xxxx-xxxx[AM][PM] OR TBA", token.Data))
				}
			}
		}
	}
	return nil
}

func getLocation(c *Course, tokenizer *html.Tokenizer, fieldCount *int,  startToken html.Token) error {
	/* getLocation gets the location and the room number 
	of the course

	Course location is what building the course is in (SLC, DDD etc.)
	Course room number is the number of the room in that building 

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/
	if (c.DeliveryMode == "SOL") {
		fmt.Printf("Course Location found: Online\n")
		output := "Online"
		c.Location = &output 
		fmt.Println("End of Location field, exiting . . .")
		return nil
	}
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Location field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Course Location found: %s\n", token.Data)
			if (token.Data == "TBA") {
				c.Location = &token.Data
			} else {
				splitData := strings.Split(token.Data, " ");
				if (len(splitData) != 2) {
					return errors.New(fmt.Sprintf("Course Location in unexpected format." + 
									   " Got %d; Expected 2.", len(splitData)))
				}
				c.Location = &splitData[0]
				roomNum, err := strconv.Atoi(splitData[1])
				if err != nil {
					c.Location = &token.Data
				} else {
					c.RoomNum = &roomNum
				}
			}
		}
	}
	return nil
}

func getInstructor(c *Course, tokenizer *html.Tokenizer, fieldCount *int,  startToken html.Token) error {
	/* getInstructor gets the instructor from the course.

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Instructor field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Instructor Token found: %s\n", token.Data)
			c.Instructor = token.Data
		}
	}
	return nil
}

func getStatus(c *Course, tokenizer *html.Tokenizer, fieldCount *int,  startToken html.Token) error {
	/* getStatus gets the status from the course.

	Course status is how full the course is (Nearly, Closed, Open)

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	textTokenCount := 0
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Status field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken && textTokenCount == 0) {

			fmt.Printf("Status Token found: %s\n", token.Data)
			c.Status = token.Data
			textTokenCount++

		} else if (tokenType == html.TextToken && textTokenCount == 1) {

			fmt.Printf("Limit Token found: %s\n", token.Data)
			limit, err := strconv.Atoi(token.Data)
			if err != nil {
				return err
			} else {
				c.Limit = limit
			}
		}
	}
	return nil
}

func getStudents(c *Course, tokenizer *html.Tokenizer, fieldCount *int,  startToken html.Token) error {
	/* getStudents gets the number of students in the course.

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Students field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Student Token found: %s\n", token.Data)
			parsedStudents, err := strconv.Atoi(token.Data)
			if err != nil {
				return err
			}
			c.Students = parsedStudents
		}
	}
	return nil
}

func getWaiting(c *Course, tokenizer *html.Tokenizer, fieldCount *int,  startToken html.Token) error {
	/* getStudents gets the number of students waiting in the course.

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Waiting field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Waiting Token found: %s\n", token.Data)
			parsedWaiting, err := strconv.Atoi(token.Data)
			if err != nil {
				return err
			}
			c.Waiting = parsedWaiting 
		}
	}
	return nil
}

func getInfo(c *Course, tokenizer *html.Tokenizer, fieldCount *int, startToken html.Token) error {
	/* getInfo gets info related to a previous course.

	Info will look something like: HONORS STUDENTS ONLY

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
		fieldCount (*int): The current field the parser is on.
		startToken (html.Token): The current token this field is starting on.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		fmt.Printf("Info: %s & %s \n", token.Data, tokenType)
		if ((tokenType == html.ErrorToken) || ((tokenType == html.EndTagToken) && (token.Data == "td")))  {
			fmt.Println("End of Info field, exiting . . .")
			break
		}

		if (tokenType == html.TextToken) {
			fmt.Printf("Info Token found: %s\n", token.Data)
			info := token.Data
			c.Info = &info
		}
	}
	return nil
}

func getField(c *Course, fieldCount *int, tokenizer *html.Tokenizer, startToken html.Token) error {
	/* getField gets the corresponding field in a 
	list of field functions from the given count

	Arguments:
		c (*course): The course to get the field for 
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
	
	Returns:
		error: Error during parsing or Error because field count is out of range
		       from avaiable functions or nil 
	*/

	fieldFuncs := []fieldFunc{
		getDeliveryMode,
		getCourseCategoryAndId,
		getSection,
		getCRN,
		getTitle,
		getCredits,
		getDay,
		getTime,
		getLocation,
		getInstructor,
		getStatus,
		getStudents,
		getWaiting,
	}

	if (c.IsCourseChild) {

		var err error
		tokenizer.Next() // <td> -> </td>
		tokenizer.Next() // </td> -> <td>
		token := tokenizer.Token()
		hasColspan := false
		colspanVal := 0
		for _, attr := range token.Attr {
			if (attr.Key == "colspan") {
				hasColspan = true
				colspanVal, _ = strconv.Atoi(attr.Val)
				break
			}
		}
		/* For the case of rows that appear as:
		 * <tr><td colspan=6></td><td colspan=7></td></tr> 
		 */
		if (hasColspan && colspanVal == 7) {
			fmt.Printf("Course Child is extra info\n")
			err = getInfo(c, tokenizer, fieldCount, startToken)

		/* For the case of rows that appear as:
		 * <tr><td colspan=6></td><td>Some Data</td><td>Some Data</td><td>Some Data</td></tr> 
		 */
		} else {
			fmt.Printf("Course Child is extra time\n")
			err = getDay(c, tokenizer, fieldCount, startToken)
			if err !=  nil { return err }
			err = getTime(c, tokenizer, fieldCount, startToken)
			if err != nil { return err }
			err = getLocation(c, tokenizer, fieldCount, startToken)
			tokenizer.Next()
		}

		return err 

	} else if (*fieldCount < len(fieldFuncs)) {
		err := fieldFuncs[*fieldCount](c, tokenizer, fieldCount, startToken)
		*fieldCount++
		return err 
	} else {
		return errors.New(fmt.Sprintf("Error: fieldCount is out of bounds of avaiable functions. " + 
								  "Got %d, Functions: %d", fieldCount, len(fieldFuncs)))
	}
}

func getCourseData (tokenizer *html.Tokenizer) (Course, error) {
	/* getCourseData parses course data from the current table row.

	It will break down the course into parts and get each field based
	on a count, starting from 0.

	Arguments:
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
	
	Returns:
		course, error: the course parsed, An error that occured during parsing or nil
	*/

	c := Course{}
	fieldCount := 0
	fmt.Println("Getting Course Data . . .")
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		if (tokenType == html.ErrorToken) {
			return c, nil
		} else if (tokenType == html.EndTagToken) && (token.Data == "tr") {
			return c, nil
		}

		fmt.Printf("In Course: Token[%s] Type[%s] fieldCount[%d]\n", token.Data, tokenType, fieldCount)

		if (tokenType == html.StartTagToken) && (token.Data == "td") {
			// Check if <td> has attribute 'colspan' and the value of it is 6
			// If so, then this row is course child
			for _, attr := range token.Attr {
				if (attr.Key == "colspan" && fieldCount == 0) {
					c.IsCourseChild = true
					fmt.Println("Course Child Found")
				}
			}
			err := getField(&c, &fieldCount, tokenizer, token)
			if err != nil {
				return c, errors.New(fmt.Sprintf("Error parsing course: %s", err)) 
			}
		}
	}
}

func getHTML(link string) (string, error) {
	/* getHTML gets the HTML from a webpage.

	Arguments:
		link (string): The link to get the HTML from.
	
	Returns:
		string: The HTML from the webpage.
	*/

	resp, err := http.Get(link)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseHTML(body string, workerNum int, wg *sync.WaitGroup, verifyChunk chan<- bool, dbChan chan<- Course, ctx context.Context) []Course {
	/* parseHTML looks at a string of HTML, and tokenizes it using Golang's html tokenizer.

	Arguments:
		body (string): The body of html to parse.
		workerNum (int): The number of this worker.
		wg (*sync.WaitGroup): The Wait Group this worker is apart of.
		verifyChunk (chan<- bool): A channel to send verification if this is a good chunk.
		dbChan (chan<- Courses): Courses will be put on this channel.
		ctx (context.Context): A context that will ensure we stop doing work if an error occurs

	Returns:
		[]Courses: A list of courses retrived.
	*/
	defer wg.Done()
	tokenizer := html.NewTokenizer(strings.NewReader(body))
	courses := []Course{}
	i := -1

	for {
		select {
		// Check the context to make sure we dont do extra work
		case <-ctx.Done():
			return []Course{}
		default:
			tokenType := tokenizer.Next()

			if (tokenType == html.ErrorToken) {

				// EOF: Add the last course to the DB
				select {
				case <-ctx.Done():
					return []Course{}
				default:
					// Send the last course to the DB
					dbChan <- courses[i]
				}
				fmt.Printf("Worker[%d]: Sent final course to DB: %s\n", workerNum, courseToString(courses[i]))
				fmt.Printf("Worker[%d]: Done.\n", workerNum)
				return courses

			}

			token := tokenizer.Token()
			fmt.Printf("Worker[%d]: Token[%s] Type[%s]\n", workerNum, token.Data, tokenType)

			if !(tokenType == html.StartTagToken && token.Data == "tr") {
				continue
			}

			c, err := getCourseData(tokenizer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Worker[%d]: %s\n", workerNum, err)
				return []Course{}
			}

			// If we havent found a course yet, check if the first course is a child.
			// If it is, then we have a bad chunk.
			if i == -1 {
				if c.IsCourseChild {
					// Bad Chunk
					verifyChunk <- false
				} else {
					// Good Chunk
					verifyChunk <- true
					courses = append(courses, c)
					i++
				}
			} else {
				if c.IsCourseChild {
					// Add child to the pervious course
					n := &courses[i]
					for n.CourseChild != nil {
						n = n.CourseChild
					}
					n.CourseChild = &c
					fmt.Printf("Worker[%d]: %s\n", workerNum, courseToString(courses[i]))
				} else {
					select {
					case <-ctx.Done():
						return []Course{}
					default:
						// Send the pervious course to the DB because it has no more children
						dbChan <- courses[i] 
					}

					fmt.Printf("Worker[%d]: Sent course to DB: %s\n", workerNum, courseToString(courses[i]))
					courses = append(courses, c)
					i++
				}
			}
		}
	}

	return courses
}
func skipToFirstRow(body string) (string, error) {
	/* skipToFirstRow takes the body and finds the first row after </thead>.

	Arguments:
		body (string): The body to slice.

	Returns:
		(string, error): The sliced body, error is not nil if </thead>
		is not found, or <tr is not found after </thead>
	*/
	endOfHead := "</thead>"
	rowStart := "<tr"

	i := strings.Index(body, endOfHead)

	if (i == -1) {
		return "", errors.New("</thead> not found")
	}

	fragmentBody := body[i:]

	j := strings.Index(fragmentBody, rowStart)
	
	if (j == -1) {
		return "", errors.New("<tr> after </thead> not found")
	}
	
	return fragmentBody[j:], nil
}

func getChunks(body string, numChunks int, shifts int) ([]string, error) {
	/* getChunks divides the body into chunks based on the value of numChunks.

	The chunks are also divided base on the end of rows. Thus chunks are all not the
	same size. They are pretty close to one another, about ~1000 charatcers

	Arguments:
		body (string): The body to divide into chunks.
		numChunks (int): The number of chunks to divide the body into.
		shifts (int): The number of left row shifts to apply. Mainly used to handle
				 related row spliting.
	Returns:
		([]string, error): The sliced body. Error is not nil if an error occurs during chunking.
	*/
	rowEnd := "</tr>"
	chunks := []string{}
	chunkSize := len(body) / numChunks
	fmt.Printf("Chunk Size: %d\n", chunkSize)
	i := 0
	chunkIndex := i + chunkSize

	for {
		// If the next chunk will be out of bounds, set the last chunk to the rest of the body.
		// This means the last chunk will most likey always be the smallest chunk while the others
		// are close to the wanted chunk size.

		if (chunkIndex > len(body)) {
			chunks = append(chunks, body[i:])
			break
		}
		chunk := body[i:chunkIndex] // Get a chunk from the body
		j := strings.LastIndex(chunk, rowEnd) // Find the last </tr> in this chunk

		if (j == -1) {
			// Set this chunk to the rest of the body
			chunks = append(chunks, body[i:])
			break
		}
	
		// Apply row left shift by finding a </tr> before index j
		leftShifts := 0
		for (leftShifts < shifts) {
			j = strings.LastIndex(chunk[0:j], rowEnd)
			if (j == -1) {
				return []string{}, errors.New(fmt.Sprintf("Too little rows in chunk, failed to shift up by %d rows", shifts))
			}
			leftShifts++
		}

		end := i + j + len(rowEnd)
		chunk = body[i:end] // Get a chunk
		chunks = append(chunks, chunk)

		i = end // Skip to the first row not in the last chunk
		chunkIndex = i + chunkSize
	}

	return chunks, nil
}

var m sync.Mutex
var count int32 = 0
func inserter(coursesIn <-chan Course, group string, wg *sync.WaitGroup) {
	/* inserters put a course into the database.

	Arguments:
		coursesIn (<-chan Course): Courses sent from the parsers to put into the database.
		group (string): The group that this course is apart of.
		wg (*sync.WaitGroup): The waitgroup the inserter is apart of.
	*/

	defer wg.Done()
	for c := range coursesIn {
		atomic.AddInt32(&count, 1)
		course, _ := json.Marshal(c)
		api.InsertCourse(course, group)
		fmt.Println("Inserter put a course in a database")
	}	
}

func Scraper() {
	/* Scraper takes 2 command line arguments, and parses the
	The Wilkes Univeristy's Course Registar pages.

	usage: scraper [F | Sp] year
	*/

	log.Println("Scraper service started")

	// Get args
	args := os.Args

	if (len(args) != 2) {
		panic(errors.New("usage: scraper semester year"))
	}

	semester := args[0]
	year := args[1]

	// Verify the semester
	possibleSemesters := []string{"F","Sp"}
	found := false
	for i := range(possibleSemesters) {
		if (semester == possibleSemesters[i]) {
			found = true
		}
	}
	if !(found) {
		panic(errors.New(fmt.Sprintf("error: bad semester, Got %s", semester)))
	}

	// Verify the year
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		panic(errors.New(fmt.Sprintf("error: bad year, Got %s", year)))
	}

	body, err := getHTML(fmt.Sprintf("https://rosters.wilkes.edu/scheds/courses%s%d.html", semester, yearInt))
	if err != nil {
		panic(err)
	}

	body, err = skipToFirstRow(body)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	// 3 inserters, the rest are parsers
	inserters := 3
	parsers := 10

	shifts := 0
	sendDB := make(chan Course, parsers)
	group := semester + year
	var insertersWg sync.WaitGroup

	// Create inserters
	for _ = range(inserters) {
		insertersWg.Add(1)
		go inserter(sendDB, group, &insertersWg)
	}

	for {
		ctx, cancel := context.WithCancel(context.Background())

		// Get chunks of the body
		chunks, err := getChunks(body, parsers - 1, shifts)
		if err != nil {
			panic(err)
		} else if len(chunks) != parsers {
			panic(errors.New(fmt.Sprintf("error: chunks do not match requested parsers. chunks: %d, parsers: %d",len(chunks), parsers)))
		}

		verifyChan := make(chan bool, parsers)

		var parsersWg sync.WaitGroup
		
		// Spawn workers
		for i := range(parsers) {
			parsersWg.Add(1)
			go parseHTML(chunks[i], i, &parsersWg, verifyChan, sendDB, ctx)
		}

		// Wait for each worker to determine if there chunk is good or bad.
		// If we find out a chunk is bad, redo the chunks
		goodChunks := 0
		badChunkFound := false
		for res := range verifyChan {
			if (!res) {
				cancel()
				// Wait for the parsers to finish
				parsersWg.Wait()
				badChunkFound = true
				break
			}
			goodChunks++
			if (goodChunks == parsers) {
				close(verifyChan)
			}
		}

		// If a bad chunk is found, increase the shift and try again
		if (badChunkFound) {
			shifts++

		} else {

			parsersWg.Wait()
			break

		}
	}

	// The parsers are done, so close the channel
	close(sendDB)
	insertersWg.Wait()

	fmt.Printf("Done Scraping. Parsed %d Courses from %s%s\n", count, semester, year)
}
