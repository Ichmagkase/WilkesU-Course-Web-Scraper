package main

import (
    "fmt"
    // "net/http"
    //"time"
)

func main() {

	// make URL
	// url := "https://example.com"

	// // Make the GET request
	// resp, _ := http.Get(url)
	// defer resp.Body.Close()

	// // Read the response body
	// bytes, _ := io.ReadAll(resp.Body)

	// // Print HTML
	// fmt.Println("HTML:\n\n", string(bytes))

	// // Connect to mongodb host
	// client, ctx, cancel, err := connect("mongodb://localhost:27017")
    // if err != nil {
    //     panic(err)
    // }

    // // Release mongodb resource when the main
    // // function is returned.
    // defer close(client, ctx, cancel)

	fmt.Println("Main executed")
	scraper()
	test()
}
