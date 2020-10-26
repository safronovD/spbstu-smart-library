## Usage

```ShellSession
    make run flags="-launch-mod=download-pdf"
```

Launch mods:

- download-json : load json data to Elasticsearch or filesystem with specified request
- dowmload-pdf : load pdf files using hrefs from csv file (now after download-json)
- samples : temp mod for dev

Other options:

- -log-file=connector.log
- -output-dir=output
- -config-file=config.yaml

Login to https://elib.spbstu.ru/ , see your .ASPXAUTH and ASP.NET_SessionId cookies and add values in config.yaml" (Pure implementation, we will develop normal authorization later :)

For Firefox: F12 -> Storage -> Cookies

## Sample requests

records' list - https://ruslan.library.spbstu.ru/rrs-web/db/EBOOKS?query=(dc.type="Academic thesis")and(dc.language=rus)&startRecord=1&maximumRecords=4&fcq=(bib.dateIssued = "2018")

one record - https://ruslan.library.spbstu.ru/rrs-web/db/EBOOKS/RU%5CSPSTU%5Cedoc%5C20151?recordSchema=gost-7.0.100
