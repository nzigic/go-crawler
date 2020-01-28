package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

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
					Url: url,
				}

				if extractErr != nil {
					crawlResult.Broken = true
					crawlResult.Message = extractErr.Error()
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
		if hrefExists && strings.HasPrefix(link, "/") {
			r = append(r, link)
		}
	})

	return
}
