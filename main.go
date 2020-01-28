package main

import (
	crawler "crawler/services"
)

func main() {
	crawler.Crawl("http://bestbytes.de")
}
