package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"strconv"
	"golang.org/x/net/html"
	"errors"
	"os"
)

/* The course struct is what a course is expected to look like.

Any variable that is a pointer to a type is an optional paramteter.
*/
type course struct {
	deliveryMode string // F2F; HYB; null etc.
	courseCategory string // CS; MTH; ENG etc.
	courseId int
	section string // INA; A; B etc.
	crn int
	title string // Planet Earth; Composition; Calculus I etc.
	credits float32 // 3.00; 4.00 etc.
	day *string // MWF; TR; null etc.
	startTime *string // 0100; 0800; 0430; null etc.
	endTime *string // 0100; 0800; 0430; null etc.
	endTimeAMPM *string // AM; PM; null.
	location *string // SLC 409; BREIS 018; null etc. 
	instructor string // Nye B; Simpson H; Kapolka M etc.
	status string // Open; Nearly; Closed.
	students int 
	waiting int
	info *string // HONORS STUDENTS ONLY; CROSS-LISTED WITH IM 350 A etc.
	isOnline bool // This is for full online classes (OL) not SOL or HYB
	child *courseChild
}

/* courseChild is for more time slots when courses dont always
meet at the same time. 

For example, CS 125 may meet on MW 0100-0215PM,
but on F 0900-1050AM. This is probably a lab, but is not listed as a lab,
thus we add a courseChild
*/
type courseChild struct {
	credits float32
	day *string
	startTime *string
	endTime *string
	endTimeAMPM *string
	location *string
	info *string
	child *courseChild
}

func getDeliveryMode(c *course, tokenizer *html.Tokenizer) error {
	/* getDeliveryMode gets the delivery mode from the course.

	Delivery modes is how the course is offered (F2F, HYB, OL etc.)

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		if (tokenType == html.ErrorToken) {
			fmt.Println("Error Token, exiting . . .")
			break
		}

		token := tokenizer.Token()

		if (tokenType == html.TextToken) {
			fmt.Printf("Delivery Mode Token found: %s\n", token.Data)
			c.deliveryMode = token.Data
			break
		}
	}
	return nil
}

func getCourseCategoryAndId(c *course, tokenizer *html.Tokenizer) error {
	/* getCourseCategoryAndId gets the course category 
	and id of the course.

	Course category is where the course belongs to (CS, MTH, ENG etc.)
	Course id is the number of that course. 

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
	
	Returns:
		error: Error during parsing or nil 
	*/
	for {
		tokenType := tokenizer.Next()
		if (tokenType == html.ErrorToken) {
			fmt.Println("Error Token, exiting . . .")
			break
		}

		token := tokenizer.Token()

		if (tokenType == html.TextToken) {
			fmt.Printf("Course Category and Id found: %s\n", token.Data)
			splitData := strings.Split(token.Data, " ");
			if (len(splitData) != 2) {
				return errors.New(fmt.Sprintf("Course category and id in unexpected format." + 
									   " Got %d; Expected 2.", len(splitData)))
			}
			c.courseCategory = splitData[0]
			courseId, err := strconv.Atoi(splitData[1])
			if err != nil {
				return err
			}
			c.courseId = courseId
			break
		}
	}
	return nil
}

func getSection(c *course, tokenizer *html.Tokenizer) error {
	/* getSection gets the section from the course.

	Section is the group of the Section (A, B, INA etc.)

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		if (tokenType == html.ErrorToken) {
			fmt.Println("Error Token, exiting . . .")
			break
		}

		token := tokenizer.Token()

		if (tokenType == html.TextToken) {
			fmt.Printf("Section Token found: %s\n", token.Data)
			c.section = token.Data
			break
		}
	}
	return nil
}

func getCRN(c *course, tokenizer *html.Tokenizer) error {
	/* getCRN gets the CRN from the course.

	CRN is the course registration number (31233,23213, 0 etc.)

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		if (tokenType == html.ErrorToken) {
			fmt.Println("Error Token, exiting . . .")
			break
		}

		token := tokenizer.Token()

		if (tokenType == html.TextToken) {
			fmt.Printf("CRN Token found: %s\n", token.Data)
			parsedCRN, err := strconv.Atoi(token.Data)
			if err != nil {
				return err
			}
			c.crn = parsedCRN
			break
		}
	}
	return nil

}

func getTitle(c *course, tokenizer *html.Tokenizer) error {
	/* getTitle gets the name from the course.

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		if (tokenType == html.ErrorToken) {
			fmt.Println("Error Token, exiting . . .")
			break
		}

		token := tokenizer.Token()

		if (tokenType == html.TextToken) {
			fmt.Printf("Title Token found: %s\n", token.Data)
			c.title = token.Data
			break
		}
	}
	return nil

}

func getCredits(c *course, tokenizer *html.Tokenizer) error {
	/* getCredits gets the credits from the course.

	Credits are floats (3.00, 4.00, 1.00, etc.)

	Arguments:
		c (*course): The course to add the delivery mode to.
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
	
	Returns:
		error: Error during parsing or nil 
	*/

	for {
		tokenType := tokenizer.Next()
		if (tokenType == html.ErrorToken) {
			fmt.Println("Error Token, exiting . . .")
			break
		}

		token := tokenizer.Token()

		if (tokenType == html.TextToken) {
			fmt.Printf("Credits Token found: %s\n", token.Data)
			parsedCredits, err := strconv.ParseFloat(token.Data, 32)
			if err != nil {
				return err
			}
			c.credits = float32(parsedCredits)
			break
		}
	}
	return nil


}

func getDay(c *course, tokenizer *html.Tokenizer) error {
	return nil
}

func getTime(c *course, tokenizer *html.Tokenizer) error {
	return nil
}

func getLocation(c *course, tokenizer *html.Tokenizer) error {
	return nil
}

func getInstructor(c *course, tokenizer *html.Tokenizer) error {
	return nil
}

func getStatus(c *course, tokenizer *html.Tokenizer) error {
	return nil
}

func getWaiting(c *course, tokenizer *html.Tokenizer) error {
	return nil
}

func getInfo(c *course, tokenizer *html.Tokenizer) error {
	return nil
}

func getField(c *course, fieldCount int, tokenizer *html.Tokenizer) error {
	switch (fieldCount) {
		case 0:
			return getDeliveryMode(c, tokenizer)
		case 1:
			return getCourseCategoryAndId(c, tokenizer)
		case 2:
			return getSection(c, tokenizer)
		case 3:
			return getCRN(c, tokenizer)
		case 4:
			return getTitle(c, tokenizer)
		case 5:
			return getCredits(c, tokenizer)
		case 6:
			return getDay(c, tokenizer)
		case 7:
			return getTime(c, tokenizer)
		case 8:
			return getLocation(c, tokenizer)
		case 9:
			return getInstructor(c, tokenizer)
		case 10:
			return getStatus(c, tokenizer)
		case 11:
			return getWaiting(c, tokenizer)
	}
	return nil
}

func getCourseData (tokenizer *html.Tokenizer) (course, error) {
	/* getCourseData parses course data from the current table row.

	It will break down the course into parts and get each field based
	on a count, starting from 0.

	Arguments:
		tokenizer (*html.Tokenizer): The tokenizer to use to get the data.
	
	Returns:
		course, error: the course parsed, An error that occured during parsing or nil
	*/

	c := course{}
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

		fmt.Printf("Token: %s ; Type: %s ; fieldCount: %d\n", token.Data, tokenType, fieldCount)

		if (tokenType == html.StartTagToken) && (token.Data == "td") {
			err := getField(&c, fieldCount, tokenizer)
			fieldCount++
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing course: %s\n", err)
				return c, err
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

func parseHTML(body string) {
	/* parseHTML looks at a string of HTML, and tokenizes it using
	Golang's html tokenizer.

	Arguments:
		body (string): The body of html to parse.
	
	Returns:
		Umm . . . Cheddar!
	*/
	tokenizer := html.NewTokenizer(strings.NewReader(body))

	// skip to tbody
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return
		}
		token := tokenizer.Token()
		fmt.Printf("Token: %s ; Type: %s\n", token.Data, tokenType)
		if (tokenType == html.EndTagToken) && (token.Data == "thead") {
			fmt.Println("At <tbody>, exiting . . .")
			tokenizer.Next() // </thead> -> <tbody>
			break
		}
	}

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			// EOF: Done reading
			return
		}
		token := tokenizer.Token()
		fmt.Printf("Token: %s ; Type: %s\n", token.Data, tokenType)
		if (tokenType == html.StartTagToken) && (token.Data == "tr") {
			getCourseData(tokenizer)
		}
	}
}

func scraper() {
	fmt.Println("Scraper service started")

	body, err := getHTML("https://rosters.wilkes.edu/scheds/coursesF25.html")
	if err != nil {
		panic(err)
	}

	parseHTML(body)
}
