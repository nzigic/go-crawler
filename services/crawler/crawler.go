package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type CrawlerService struct {
}

type CrawlResult struct {
	Url     string
	Broken  bool
	Message string
}

type CrawlRequest struct {
	rootUrl string
}

func New() *CrawlerService {
	return &CrawlerService{}
}

func (s *CrawlerService) Crawl(rootUrl string) (r []CrawlResult) {
	r = []CrawlResult{}
	var wg sync.WaitGroup

	visitedUrls := make(map[string]bool)
	chanExtractedLinks := make(chan string, 5)

	wg.Add(1)
	go func() {
		chanExtractedLinks <- rootUrl
	}()

	go func() {
		defer close(chanExtractedLinks)
		wg.Wait()
	}()

	isFirstRun := true
	var linkCounter int64
	var averageElapsedPerPage float64
	start := time.Now()

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
					Url: url,
				}

				if extractErr != nil {
					crawlResult.Broken = true
					crawlResult.Message = extractErr.Error()
				} else {
					validLinks := filterInternalLinks(extractedLinks, rootUrl)
					for _, link := range validLinks {
						chanExtractedLinks <- link
					}
				}

				r = append(r, crawlResult)

				linkCounter++
				elapsed := time.Since(start)
				averageElapsedPerPage = float64(elapsed.Milliseconds()) / float64(linkCounter)
				fmt.Printf("\rProcessed: %d links, average time per URL: %fms", linkCounter, averageElapsedPerPage)

			}(extractedLink)
		}
	}

	fmt.Println()
	return
}

func extractLinks(pageUrl string) (r []string, err error) {
	resp, errGet := http.Get(pageUrl)
	if errGet != nil {
		return nil, errGet
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	doc, errNewDoc := goquery.NewDocumentFromReader(resp.Body)
	if errNewDoc != nil {
		return nil, errNewDoc
	}

	r = []string{}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, hrefExists := s.Attr("href")
		if hrefExists {
			r = append(r, link)
		}
	})

	return
}

func filterInternalLinks(extractedLinks []string, rootUrl string) []string {
	result := []string{}
	for _, link := range extractedLinks {
		if strings.HasPrefix(link, "/") {
			result = append(result, rootUrl+link)
			continue
		}

		if strings.HasPrefix(link, rootUrl) {
			result = append(result, link)
		}
	}

	return result

}
