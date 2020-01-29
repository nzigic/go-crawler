package main

import (
	"crawler/services/crawler"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	s := crawler.New()
	if len(os.Args) != 0 && (os.Args[1] == "--web" || os.Args[1] == "-w") {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		crawlerProxy := crawler.NewDefaultCrawlerServiceGoTSRPCProxy(s, []string{"*"})
		crawlerMethods := []string{"/Crawl"}
		registerMethods(mux, crawlerProxy, crawlerProxy.EndPoint, crawlerMethods)

		fmt.Println("server started on port 8080")
		http.ListenAndServe(":8080", mux)
	}

	start := time.Now()

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

	elapsed := time.Since(start)
	fmt.Printf("took %s", elapsed)
	fmt.Println()
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

func registerMethods(mux *http.ServeMux, proxy http.Handler, proxyEndpoint string, methods []string) {
	for _, method := range methods {
		mux.Handle(proxyEndpoint+method, proxy)
	}
}
