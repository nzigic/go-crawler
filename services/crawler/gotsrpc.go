// Code generated by gotsrpc https://github.com/foomo/gotsrpc  - DO NOT EDIT.

package crawler

import (
	http "net/http"
	time "time"

	gotsrpc "github.com/foomo/gotsrpc"
)

type CrawlerServiceGoTSRPCProxy struct {
	EndPoint    string
	allowOrigin []string
	service     *CrawlerService
}

func NewDefaultCrawlerServiceGoTSRPCProxy(service *CrawlerService, allowOrigin []string) *CrawlerServiceGoTSRPCProxy {
	return &CrawlerServiceGoTSRPCProxy{
		EndPoint:    "/services/crawler",
		allowOrigin: allowOrigin,
		service:     service,
	}
}

func NewCrawlerServiceGoTSRPCProxy(service *CrawlerService, endpoint string, allowOrigin []string) *CrawlerServiceGoTSRPCProxy {
	return &CrawlerServiceGoTSRPCProxy{
		EndPoint:    endpoint,
		allowOrigin: allowOrigin,
		service:     service,
	}
}

// ServeHTTP exposes your service
func (p *CrawlerServiceGoTSRPCProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	for _, origin := range p.allowOrigin {
		// todo we have to compare this with the referer ... and only send one
		w.Header().Add("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method != http.MethodPost {
		if r.Method == http.MethodOptions {
			return
		}
		gotsrpc.ErrorMethodNotAllowed(w)
		return
	}

	var args []interface{}
	funcName := gotsrpc.GetCalledFunc(r, p.EndPoint)
	callStats := gotsrpc.GetStatsForRequest(r)
	if callStats != nil {
		callStats.Func = funcName
		callStats.Package = "crawler/services/crawler"
		callStats.Service = "CrawlerService"
	}
	switch funcName {
	case "Crawl":
		var (
			arg_rootUrl string
		)
		args = []interface{}{&arg_rootUrl}
		err := gotsrpc.LoadArgs(&args, callStats, r)
		if err != nil {
			gotsrpc.ErrorCouldNotLoadArgs(w)
			return
		}
		executionStart := time.Now()
		crawlR := p.service.Crawl(arg_rootUrl)
		if callStats != nil {
			callStats.Execution = time.Now().Sub(executionStart)
		}
		gotsrpc.Reply([]interface{}{crawlR}, callStats, r, w)
		return
	default:
		gotsrpc.ClearStats(r)
		http.Error(w, "404 - not found "+r.URL.Path, http.StatusNotFound)
	}
}
