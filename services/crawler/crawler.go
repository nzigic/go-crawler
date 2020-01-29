package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CrawlerService struct {
}

type CrawlResult struct {
	Url     string `json:"url,omitempty"`
	Broken  bool   `json:"broken,omitempty"`
	Message string `json:"message,omitempty"`
}

var (
	urlsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "crawler_processed_urls_total",
		Help: "The total number of processed URLs",
	})
	successfulUrls = promauto.NewCounter(prometheus.CounterOpts{
		Name: "crawler_successful_urls_total",
		Help: "The total number of successful (HTTP 200) URLs",
	})
	failedUrls = promauto.NewCounter(prometheus.CounterOpts{
		Name: "crawler_failed_urls_total",
		Help: "The total number of failed (HTTP 500) URLs",
	})
)

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

	// todo: Add average response time counter or something similar
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

				urlsProcessed.Inc()
			}(extractedLink)
		}
	}

	fmt.Println()
	return
}

func extractLinks(pageUrl string) (r []string, err error) {
	resp, errGet := http.Get(pageUrl)
	if errGet != nil {
		failedUrls.Inc()
		return nil, errGet
	}

	if resp.StatusCode != 200 {
		failedUrls.Inc()
		return nil, errors.New(resp.Status)
	}

	successfulUrls.Inc()

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
