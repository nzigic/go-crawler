package main

import (
	crawler "crawler/services/crawler"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
)

func main() {
	s := crawler.New()
	fmt.Println(os.Args[0])
	if len(os.Args) != 0 && (os.Args[1] == "--web" || os.Args[1] == "-w") {
		proxy := crawler.NewDefaultCrawlerServiceGoTSRPCProxy(s, []string{"*"})
		fmt.Println("server started on port 8080")

		http.ListenAndServe(":8080", proxy)
	}

	fmt.Println("CLI running...")
	results := s.Crawl("http://bestbytes.de")
	brokenUrls := filterBrokenLinks(results)
	fmt.Println("Total URLs: ", len(results))
	fmt.Println("Broken URLs: ")
	for _, link := range brokenUrls {
		color.Set(color.FgHiRed, color.Bold)
		defer color.Unset()
		fmt.Printf("%s: '%s'", link.Url, link.Message)
		fmt.Println()
	}
}

func filterBrokenLinks(links []crawler.CrawlResult) (out []crawler.CrawlResult) {
	result := make([]crawler.CrawlResult, 0)
	for _, link := range links {
		if link.Broken {
			result = append(result, link)
		}
	}

	return result
}
