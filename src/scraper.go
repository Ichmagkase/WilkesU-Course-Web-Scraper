package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"golang.org/x/net/html"
)


func getHTML(link string) string {
	/* getHTML gets the HTML from a webpage.

	Arguments:
		link (string): The link to get the HTML from.
	
	Returns:
		string: The HTML from the webpage.
	*/

	resp, err := http.Get(link)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
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
		fmt.Println(token.String())
	}
}
func scraper() {
	fmt.Println("Scraper service started")

	body := getHTML("https://rosters.wilkes.edu/scheds/coursesF25.html")

	parseHTML(body)
}
