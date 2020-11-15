# Usage

## Dev

!! In recordId "\\" is replaced with "_", cause it's not a valid file name (Example: ru\spstu\vkr\1000 -> ru_spstu_vkr_1000)

output/*.csv contains list of IDs and links to pdf (for py scripts)

getting href for pdf downloading isn't working with converted jsons

## Using Docker for building

Go to build/ and run script build.sh

Use arg "remote" for building from gitHub

Error log -> log.err, output -> build/output


## Using Make with source files
```ShellSession
make run flags="-launch-mod=download-pdf"
```

## Launch mods:

- download-json : load json data to Elasticsearch or filesystem with specified request
- dowmload-pdf : load pdf files using hrefs from csv file (now after download-json)
- samples : temp mod for dev

## Other options:

- -log-file=connector.log
- -output-dir=output
- -config-file=config.yaml

Login to https://elib.spbstu.ru/ , see your .ASPXAUTH and ASP.NET_SessionId cookies and add values in config.yaml" (Pure implementation, we will develop normal authorization later :)

For Firefox: F12 -> Storage -> Cookies

## Elasticsearch setup

```ShellSession
sudo docker run -d --rm \
    -p 9200:9200 -p 9300:9300 \
    -e ELASTIC_USERNAME="" -e ELASTIC_PASSWORD="" \
    -e "discovery.type=single-node" \
    -e "xpack.security.enabled=true" \
    --name elasticsearch \
    docker.elastic.co/elasticsearch/elasticsearch:7.9.3

sudo docker run -d --rm \
    -p 5601:5601 \
    -e ELASTICSEARCH_HOSTS="http://192.168.116.215:9200" \
    -e ELASTICSEARCH_USERNAME="" -e ELASTICSEARCH_PASSWORD="" \
    -e "xpack.security.enabled=true" \
    --name kibana \
    docker.elastic.co/kibana/kibana:7.9.3
```

## Sample requests

records' list - https://ruslan.library.spbstu.ru/rrs-web/db/EBOOKS?query=(dc.type="Academic thesis")and(dc.language=rus)&startRecord=1&maximumRecords=4&fcq=(bib.dateIssued = "2018")

one record - https://ruslan.library.spbstu.ru/rrs-web/db/EBOOKS/RU%5CSPSTU%5Cedoc%5C20151?recordSchema=gost-7.0.100
