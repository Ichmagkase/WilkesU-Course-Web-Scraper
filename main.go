package main

import (
    "fmt"
	"sync"
	"wilkesu-scrapy/scraper"
)

func main() {
	fmt.Println("Main executed")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scraper.DatabaseIntializer()
	}()

	scraper.Scraper()
	wg.Wait()
	// scraper.ExampleInsertion()
}
