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
	chanExtractedLinks := make(chan string)

	wg.Add(1)
	go func() {
		chanExtractedLinks <- rootUrl
	}()

	go func() {
		defer close(chanExtractedLinks)
		wg.Wait()
	}()

	isFirstRun := true
	for extractedLink := range chanExtractedLinks {
		if !visitedUrls[extractedLink] {
			visitedUrls[extractedLink] = true
			wg.Add(1)
			if isFirstRun {
				isFirstRun = false
				wg.Done()
			}

			go func(url string) {
				defer wg.Done()
				extractedLinks, extractErr := extractLinks(url)
				crawlResult := CrawlResult{
					url: url,
				}

				if extractErr != nil {
					crawlResult.broken = true
				} else {
					for _, relativeUrl := range extractedLinks {
						newUrl := rootUrl + relativeUrl
						chanExtractedLinks <- newUrl
					}
				}

				r = append(r, crawlResult)
			}(extractedLink)
		}
	}

	fmt.Println(rootUrl, "processed")
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
