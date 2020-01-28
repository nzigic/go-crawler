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
	broken bool
}

func Crawl(rootUrl string) (r []CrawlResult, err error) {
	r = []CrawlResult{}
	var processWg sync.WaitGroup
	var readWg sync.WaitGroup

	visitedUrls := make(map[string]bool)
	processedUrls := make(map[string]bool)
	chanExtractedLinks := make(chan string)
	chanCrawlResults := make(chan CrawlResult)

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

	go func() {
		for extractedLink := range chanExtractedLinks {
			if !visitedUrls[extractedLink] {
				visitedUrls[extractedLink] = true
				processWg.Add(1)

				go func(url string) {
					defer processWg.Done()
					extractedLinks, extractErr := extractLinks(url)
					if extractErr != nil {
						chanCrawlResults <- CrawlResult{
							url:    url,
							broken: true,
						}
					}

					for _, relativeUrl := range extractedLinks {
						newUrl := rootUrl + relativeUrl
						chanExtractedLinks <- newUrl
					}
				}(extractedLink)
			}

			if !processedUrls[extractedLink] {
				processedUrls[extractedLink] = true

				processWg.Add(1)
				go func(url string) {
					defer processWg.Done()

					crawlResult := processLink(url)
					chanCrawlResults <- crawlResult
					readWg.Add(1)
				}(extractedLink)
			}
		}
	}()

	for crawlResult := range chanCrawlResults {
		readWg.Add(1)
		fmt.Println(crawlResult)
		r = append(r, crawlResult)
		defer readWg.Done()

		fmt.Println(&readWg)
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
	_, errGet := http.Get(link)
	if errGet != nil {
		return CrawlResult{
			url:    link,
			broken: true,
		}
	}

	r = CrawlResult{
		url: link,
	}

	return
}
