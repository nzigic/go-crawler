package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type CrawlResult struct {
	url    string
	status int
}

func Crawl(rootUrl string) (r []CrawlResult, err error) {
	r = []CrawlResult{}
	var processWg sync.WaitGroup
	var readWg sync.WaitGroup

	visitedUrls := make(map[string]bool)
	processedUrls := make(map[string]bool)
	chanExtractedLinks := make(chan string)
	chanCrawlResults := make(chan CrawlResult)

	processWg.Add(1)
	readWg.Add(1)
	go func() {
		chanExtractedLinks <- rootUrl
	}()

	go func() {
		defer close(chanExtractedLinks)
		processWg.Wait()
	}()

	go func() {
		defer close(chanCrawlResults)
		readWg.Wait()
	}()

	for extractedLink := range chanExtractedLinks {
		if !visitedUrls[extractedLink] {
			visitedUrls[extractedLink] = true
			processWg.Add(1)

			go func(l string) {
				processWg.Done()
				extractedLinks, extractErr := extractLinks(extractedLink)
				if extractErr != nil {
					chanCrawlResults <- CrawlResult{
						url:    l,
						status: 500,
					}
				}

				for i, _ := range extractedLinks {
					url := rootUrl + extractedLinks[i]
					chanExtractedLinks <- url
				}
			}(extractedLink)
		}

		if !processedUrls[extractedLink] {
			processedUrls[extractedLink] = true

			processWg.Add(1)
			go func(url string) {
				processWg.Done()

				crawlResult := processLink(url)
				chanCrawlResults <- crawlResult
			}(extractedLink)
		}

		fmt.Println("loop ", &processWg)
	}

	for crawlResult := range chanCrawlResults {
		readWg.Add(1)
		fmt.Println(crawlResult)
		r = append(r, crawlResult)
		readWg.Done()
	}

	fmt.Println("end")
	return
}

func extractLinks(pageUrl string) (r []string, err error) {
	resp, errGet := http.Get(pageUrl)
	if errGet != nil {
		return nil, errGet
	}

	defer resp.Body.Close()
	doc, errNewDoc := goquery.NewDocumentFromReader(resp.Body)
	if errNewDoc != nil {
		return nil, errNewDoc
	}

	r = []string{}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, hrefExists := s.Attr("href")
		if hrefExists && strings.HasPrefix(link, "/") {
			r = append(r, link)
		}
	})

	return
}

func processLink(link string) (r CrawlResult) {
	resp, errGet := http.Get(link)
	if errGet != nil {
		fmt.Println(" ERROR GET: ", errGet)
		return CrawlResult{
			url:    link,
			status: 500,
		}
	}

	r = CrawlResult{
		url:    link,
		status: resp.StatusCode,
	}

	return
}
