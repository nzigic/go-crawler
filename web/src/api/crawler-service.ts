/* tslint:disable */
// hello commonjs - we need some imports - sorted in alphabetical order, by go package
import * as crawler_services_crawler from './crawler-vo'; // web/src/api/crawler-service.ts to web/src/api/crawler-vo.ts

export class CrawlerServiceClient {
	public static defaultEndpoint = "/services/crawler";
	constructor(
		public transport:<T>(method: string, data?: any[]) => Promise<T>
	) {}
	async crawl(rootUrl:string):Promise<crawler_services_crawler.CrawlResult[]> {
		return (await this.transport<{0:crawler_services_crawler.CrawlResult[]}>("Crawl", [rootUrl]))[0]
	}
}