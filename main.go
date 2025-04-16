package main

import (
	"fmt"
	"os"
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

	os.Args = []string{"F","25"}
	scraper.Scraper()
	wg.Wait()
	// scraper.ExampleInsertion()
}
