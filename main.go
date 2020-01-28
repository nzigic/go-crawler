package main

import (
	crawler "crawler/services/crawler"
	"fmt"
	"net/http"
)

func main() {
	s := crawler.New()

	proxy := crawler.NewDefaultCrawlerServiceGoTSRPCProxy(s, []string{"*"})
	fmt.Println("server started on port 8080")

	http.ListenAndServe(":8080", proxy)

	// results, err := s.Crawl("http://bestbytes.de")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// brokenUrls := filterBrokenLinks(results)
	// fmt.Println("Total URLs: ", len(results))
	// fmt.Println("Broken URLs: ")
	// for _, link := range brokenUrls {
	// 	color.Set(color.FgHiRed, color.Bold)
	// 	defer color.Unset()
	// 	fmt.Printf("%s: '%s'", link.Url, link.Message)
	// 	fmt.Println()
	// }
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
