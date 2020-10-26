package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func downloadSamples(outputDir string) {
	connectUrl := "https://ruslan.library.spbstu.ru/rrs-web/db/"
	httpClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	commonPath := path.Join(".", outputDir, "samples")
	os.Mkdir(commonPath, os.ModePerm)

	req, err := http.NewRequest("PROPFIND", connectUrl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")

	data, err := downloadJson(&httpClient, req)
	if err != nil {
		log.Panic(err)
	}

	saveJson(data, path.Join(commonPath, "db_list"+".json"))

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	dbList := result["response"].([]interface{})

	for _, db := range dbList {
		href := db.(map[string]interface{})["href"].([]interface{})[0].(string)
		splitedHref := strings.Split(href, "/")
		dbName := splitedHref[len(splitedHref)-1]

		dbPath := path.Join(commonPath, dbName)
		os.Mkdir(dbPath, os.ModePerm)

		if dbName == "db" {
			continue
		}

		req, err := http.NewRequest(http.MethodGet, href, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Accept", "application/json")

		data, err := downloadJson(&httpClient, req)
		if err != nil {
			log.Panic(err)
		}

		saveJson(data, path.Join(dbPath, "db_index"+".json"))

		req, err = http.NewRequest(http.MethodGet, href, nil)
		if err != nil {
			log.Panic(err)
		}

		req.Header.Set("Accept", "application/json")

		q := req.URL.Query()
		q.Add("query", "(bib.volume=*)")
		q.Add("maximumRecords", "2")
		q.Add("startRecord", "1")
		req.URL.RawQuery = q.Encode()

		data, err = downloadJson(&httpClient, req)
		if err != nil {
			log.Panic(err)
		}

		saveJson(data, path.Join(dbPath, "db_records_list"+".json"))

		rl := RecordsList{}

		err = json.Unmarshal(data, &rl)
		if err != nil {
			log.Panic(err)
		}

		log.Printf("DB: %s Number of records: %d\n", dbName, rl.NumberOfRecords)

		for i, val := range rl.InnerRecordsList.Records {
			downloadUrl := href + "/" + url.PathEscape(val.RecordIdentifier)

			req, err := http.NewRequest(http.MethodGet, downloadUrl, nil)
			if err != nil {
				log.Panic(err)
			}

			req.Header.Set("Accept", "application/json")

			q := req.URL.Query()
			q.Add("recordSchema", "gost-7.0.100")
			req.URL.RawQuery = q.Encode()

			jsonData, err := downloadJson(&httpClient, req)
			if err != nil {
				log.Println(err)
				continue
			}

			saveJson(jsonData, path.Join(dbPath, "record"+strconv.Itoa(i)+".json"))
		}
	}
}
