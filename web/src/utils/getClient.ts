import Axios, { AxiosResponse } from "axios";
import { CrawlerServiceClient } from "../api/crawler-service";

const getTransport = (endpoint: string) => async function <T>(method: string, args: any[] | undefined) {
    return new Promise<T>(async (resolve, reject) => {
        try {
            let axiosPromise: AxiosResponse<T> = await Axios.post<T>(
                endpoint + "/" + encodeURIComponent(method),
                JSON.stringify(args),
            );
            return resolve(axiosPromise.data);
        } catch (e) {
            return reject(e);
        }
    });
};

export default function getClient() {
    return new CrawlerServiceClient(getTransport('http://localhost:8080' + CrawlerServiceClient.defaultEndpoint))
}