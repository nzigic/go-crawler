/* tslint:disable */
// hello commonjs - we need some imports - sorted in alphabetical order, by go package
import * as crawler_services_crawler from './crawler-vo'; // web/src/api/crawler-vo.ts to web/src/api/crawler-vo.ts
// crawler/services/crawler.CrawlResult
export interface CrawlResult {
	url?:string;
	broken?:boolean;
	message?:string;
}
// end of common js