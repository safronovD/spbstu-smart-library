package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	RecordIdentifier string `json:"recordIdentifier"`
}

type RecordsList struct {
	Records []Record `json:"record"`
}

type RecordsSuperList struct {
	RList              RecordsList `json:"records"`
	NumberOfRecords    int         `json:"numberOfRecords"`
	NextRecordPosition int         `json:"nextRecordPosition"`
}

func main() {

	url := "https://ruslan.library.spbstu.ru/rrs-web/db/EBOOKS"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	f, err := os.Create("dir/ids.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n := 1

	for n < 20000 {

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Accept", "application/json")

		q := req.URL.Query()
		q.Add("query", "(dc.type=\"Academic thesis\") and (dc.language=rus)")
		q.Add("fcq", "(bib.dateIssued = \"2018\")")
		q.Add("maximumRecords", strconv.Itoa(10))
		q.Add("startRecord", strconv.Itoa(n))

		req.URL.RawQuery = q.Encode()

		res, getErr := spaceClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		data, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		rl := RecordsSuperList{}

		err = json.Unmarshal(data, &rl)
		if err != nil {
			panic(err)
		}

		fmt.Println(rl.NextRecordPosition)

		for _, val := range rl.RList.Records {
			f.Write([]byte(strconv.Itoa(n)))
			f.Write([]byte(","))
			f.Write([]byte(val.RecordIdentifier + "\n"))

			jsonFile, err := os.Create("dir/jsons/" + strconv.Itoa(n) + ".json")
			if err != nil {
				panic(err)
			}
			defer jsonFile.Close()

			newURL := url + "/" + strings.ReplaceAll(val.RecordIdentifier, "\\", "%5C")

			req, err := http.NewRequest(http.MethodGet, newURL, nil)
			if err != nil {
				panic(err)
			}

			req.Header.Set("Accept", "application/json")

			q := req.URL.Query()
			q.Add("recordSchema", "gost-7.0.100")

			req.URL.RawQuery = q.Encode()

			res, getErr := spaceClient.Do(req)
			if getErr != nil {
				log.Fatal(getErr)
			}

			if res.Body != nil {
				defer res.Body.Close()
			}

			data, readErr := ioutil.ReadAll(res.Body)
			if readErr != nil {
				log.Fatal(readErr)
			}

			jsonFile.Write(data)

			rl := RecordsSuperList{}

			err = json.Unmarshal(data, &rl)
			if err != nil {
				panic(err)
			}

			n += 1
		}
	}
}
