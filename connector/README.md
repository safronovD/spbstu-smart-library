## Dev

downloadRecords or downloadSamples are hardcoded and not in config. Sorry :)

## Usage

Login to https://elib.spbstu.ru/ , see your .ASPXAUTH and ASP.NET_SessionId cookies and add values in config.yaml" (Pure implementation, we will develop normal authorization later :)

```ShellSession
    make run
```

## Sample requests

records' list - https://ruslan.library.spbstu.ru/rrs-web/db/EBOOKS?query=(dc.type="Academic thesis")and(dc.language=rus)&startRecord=1&maximumRecords=4&fcq=(bib.dateIssued = "2018")

one record - https://ruslan.library.spbstu.ru/rrs-web/db/EBOOKS/RU%5CSPSTU%5Cedoc%5C20151?recordSchema=gost-7.0.100
