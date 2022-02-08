# budget-google

Simple Google-like service

![image](https://user-images.githubusercontent.com/18232940/153023044-e37b29e8-4070-42c8-b213-df0b9f6612a6.png)

# Services:
 <h3>APIUrlToIndex</h3>

    Is running on port <:9001> Takes url and sends it to Crawler via kafka

<h3>Crawler</h3>

    Takes Url, extracts text from html page and posts it to ElasticSearch Database
 
<h3>APISearchRequest</h3>

    Is running on port <:9002> Allows user to search websites by keywords. Almost just like Google does (:
 
<h3>Prometheus</h3>

    Is running on port <:9090> Collects metrics from APIUrlToIndex, Crawler, APISearchRequest services

<h3>Grafana</h3>

    Is running on port <:3030> Displays all metrics in a web

