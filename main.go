package main

import (
	crawler "crawler/services"
	"fmt"
)

func main() {
	results, err := crawler.Crawl("http://bestbytes.de")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(results)
}
