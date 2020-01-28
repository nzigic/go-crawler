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
	var wg sync.WaitGroup

	visitedUrls := make(map[string]bool)
	processedUrls := make(map[string]bool)
	chanExtractedLinks := make(chan string)

	wg.Add(1)
	go func() {
		chanExtractedLinks <- rootUrl
	}()

	go func() {
		defer close(chanExtractedLinks)
		wg.Wait()
	}()

	for extractedLink := range chanExtractedLinks {
		if !visitedUrls[extractedLink] {
			visitedUrls[extractedLink] = true
			wg.Add(1)

			go func(url string) {
				wg.Done()
				extractedLinks, extractErr := extractLinks(url)
				if extractErr != nil {
					crawlResult := CrawlResult{
						url:    url,
						broken: true,
					}
					r = append(r, crawlResult)
				}

				for _, relativeUrl := range extractedLinks {
					newUrl := rootUrl + relativeUrl
					chanExtractedLinks <- newUrl
				}
			}(extractedLink)
		}

		if !processedUrls[extractedLink] {
			processedUrls[extractedLink] = true

			wg.Add(1)
			go func(url string) {
				wg.Done()

				crawlResult := processLink(url)
				r = append(r, crawlResult)
			}(extractedLink)
		}

		fmt.Println("loop", &wg)
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
