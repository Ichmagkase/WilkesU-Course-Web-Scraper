 package main

import (
	"os"
	"sync"
	"wilkesu-scrapy/api"
	"wilkesu-scrapy/scraper"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// Serve on :8080
	go func() {
		defer wg.Done()
		api.Serve()
	}()

	// if os.Args[0] == "-s" {
	os.Args = []string{"F", "25"}
	scraper.Scraper()
	// }

	// Wait for process to finish
	wg.Wait()
}
