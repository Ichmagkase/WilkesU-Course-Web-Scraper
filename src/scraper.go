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

func getCourseData (token html.Token, tokenizer *html.Tokenizer) course {
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			// EOF: Done reading
			return course {}
		}
		token := tokenizer.Token()
		if (tokenType == html.TextToken) {
			fmt.Printf("%s", token.Data)
		} else if (tokenType == html.EndTagToken) && (token.Data == "tr") {
			fmt.Println()
			return course {}
		}
	}
	return course {}
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
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			// EOF: Done reading
			return
		}
		token := tokenizer.Token()
		if (tokenType == html.StartTagToken) && ("tr" == token.Data) {
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
