package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"golang.org/x/net/html"
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

func getDeliveryMode(c *course, tokenizer *html.Tokenizer) {
	for {
		tokenType := tokenizer.Next()
		if (tokenType == html.ErrorToken) {
			fmt.Println("Error Token, exiting . . .")
			break
		}

		token := tokenizer.Token()

		if (tokenType == html.TextToken) {
			fmt.Printf("DeliveryMode Token found! %s\n", token.Data)
			c.deliveryMode = token.Data
			break
		}
	}
}

func getCourseCategoryAndId(c *course, tokenizer *html.Tokenizer) {
}

func getSection(c *course, tokenizer *html.Tokenizer) {
}

func getCRN(c *course, tokenizer *html.Tokenizer) {
}

func getTitle(c *course, tokenizer *html.Tokenizer) {
}

func getCredits(c *course, tokenizer *html.Tokenizer) {
}

func getDay(c *course, tokenizer *html.Tokenizer) {
}

func getTime(c *course, tokenizer *html.Tokenizer) {
}

func getLocation(c *course, tokenizer *html.Tokenizer) {
}

func getInstructor(c *course, tokenizer *html.Tokenizer) {
}

func getStatus(c *course, tokenizer *html.Tokenizer) {
}

func getWaiting(c *course, tokenizer *html.Tokenizer) {
}

func getInfo(c *course, tokenizer *html.Tokenizer) {
}

func getField(c *course, fieldCount *int, tokenizer *html.Tokenizer) {
	switch (*fieldCount) {
		case 0:
			getDeliveryMode(c, tokenizer)
			break
		case 1:
			getCourseCategoryAndId(c, tokenizer)
			break
		case 2:
			getSection(c, tokenizer)
			break
		case 3:
			getCRN(c, tokenizer)
			break
		case 4:
			getTitle(c, tokenizer)
			break
		case 5:
			getCredits(c, tokenizer)
			break
		case 6:
			getDay(c, tokenizer)
			break
		case 7:
			getTime(c, tokenizer)
			break
		case 8:
			getLocation(c, tokenizer)
			break
		case 9:
			getInstructor(c, tokenizer)
			break
		case 10:
			getStatus(c, tokenizer)
			break
		case 11:
			getWaiting(c, tokenizer)
			break
	}
	(*fieldCount)++
}

func getCourseData (token html.Token, tokenizer *html.Tokenizer) course {
	c := course{}
	fieldCount := 0
	fmt.Println("Getting Course Data . . .")
	for {
		tokenType := tokenizer.Next()
		if (tokenType == html.ErrorToken) {
			return c
		} else if (tokenType == html.EndTagToken) && (token.Data == "tr") {
			return c
		}

		token := tokenizer.Token()
		fmt.Printf("Token: %s ; Type: %s\n", token.Data, tokenType)

		if (tokenType == html.StartTagToken) && (token.Data == "td") {
			getField(&c, &fieldCount, tokenizer)
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
			getCourseData(token, tokenizer)
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
