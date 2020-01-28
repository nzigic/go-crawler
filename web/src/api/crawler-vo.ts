/* tslint:disable */
// hello commonjs - we need some imports - sorted in alphabetical order, by go package
import * as crawler_services_crawler from './crawler-vo'; // web/api/crawler-vo.ts to web/api/crawler-vo.ts
// crawler/services/crawler.CrawlResult
export interface CrawlResult {
	Url:string;
	Broken:boolean;
	Message:string;
}
// end of common js