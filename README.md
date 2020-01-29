# GoLang Crawler

## Crawling
Example usage:
curl localhost:8080/services/crawler/Crawl -X POST -d '["http://bestbytes.de"]'

## Setup
Add <YOUR_IP> to hosts file as server
so, the entry should be:
<YOUR_IP> server

This is necessary for Prometheus to work correctly

## Grafana
After running docker-compose up, run docker exec -it cf80857add22 grafana-cli admin reset-admin-password admin
to get <CONTAINER_ID>, run docker container list
Add a new data source, URL should be http://server:9090