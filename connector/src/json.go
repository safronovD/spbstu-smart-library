package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"exec"
	"github.com/elastic/go-elasticsearch/v8"
)

type Record struct {
	RecordIdentifier string `json:"recordIdentifier"`
}

type RecordsList struct {
	InnerRecordsList struct {
		Records []Record `json:"record"`
	} `json:"records"`

	NumberOfRecords    int `json:"numberOfRecords"`
	NextRecordPosition int `json:"nextRecordPosition"`
}

func saveJson(data []byte, path string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Failed to save json")
		}
	}()

	jsonFile, err := os.Create(path)
	if err != nil {
		log.Panic(err)
	}
	defer jsonFile.Close()

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, data, "", "    ")
	_, err = prettyJSON.WriteTo(jsonFile)
	if err != nil {
		os.Remove(jsonFile.Name())
		log.Panic(err)
	} else {
		log.Println(fmt.Sprintf("Json file \"%s\" saved", jsonFile.Name()))
	}
}

func downloadJson(client *http.Client, req *http.Request) ([]byte, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != http.StatusMultiStatus) {
		return nil, errors.New(fmt.Sprintf("Response failed: %s Status code: %d", res.Request.URL, res.StatusCode))
	}

	if res.Body != nil {
		defer res.Body.Close()
	} else {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func downloadRecords(config *JsonConfig, outputDir string) {
	saveToES := func(jsonData []byte, recordId string) {}
	saveToFS := func(jsonData []byte, number int, recordId string) {}

	if config.Output.Elasticsearch.Enable {
		cfg := elasticsearch.Config{
			Addresses: []string{
				config.Output.Elasticsearch.Host,
			},
			Username: config.Output.Elasticsearch.Login,
			Password: config.Output.Elasticsearch.Pwd,
		}

		es, err := elasticsearch.NewClient(cfg)
		if err != nil {
			log.Panic(err)
		}

		log.Println(elasticsearch.Version)
		log.Println(es.Info())

		saveToES = func(jsonData []byte, recordId string) {
			defer func() {
				if err := recover(); err != nil {
					log.Println("Failed to load data to ES")
				}
			}()

			//Poor logic, saving temporart to file to convert json. Change it later
			file, err := os.OpenFile("./tmp.json", os.O_RDWR|os.O_CREATE, os.ModePerm) 
			if err != nil {
				log.Panic(err)
			}
			defer file.Close()
			file.Write(jsonData)

			cmd := exec.Command("../../lib/utils/json_convertor.py", "./tmp.json", "./res.json")
			cmd.Run()
			
			convertedJson, err := ioutil.ReadFile("./res.json")
			if err != nil {
				log.Panic(err)
			}

			jsonData = convertedJson

			rsp, err := es.Index(config.Output.Elasticsearch.Index, bytes.NewReader(jsonData))
			if err != nil {
				log.Panic(err)
			}
			if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
				log.Panic(errors.New("Failed to load data. " + " Code: " + rsp.Status()))
			}

			log.Printf("Record with id \"%s\" send to ES.\n", recordId)
		}
	}

	if config.Output.FileSystem.Enable {
		commonPath := path.Join(".", outputDir, config.Connection.DB)
		jsonPath := path.Join(commonPath, config.Output.FileSystem.JsonDir)

		os.Mkdir(commonPath, os.ModePerm)
		os.Mkdir(jsonPath, os.ModePerm)

		csvFile, err := os.Create(path.Join(commonPath, config.Output.FileSystem.CsvFile))
		if err != nil {
			log.Panic(err)
		}
		defer csvFile.Close()

		csvWriter := csv.NewWriter(csvFile)

		saveToFS = func(jsonData []byte, number int, recordId string) {
			splitedRId := strings.Split(recordId, "\\")
			simpleId := splitedRId[len(splitedRId)-1]
			csvLine := []string{strconv.Itoa(number), recordId, simpleId, getHref(jsonData)}

			csvWriter.Write(csvLine)
			csvWriter.Flush()
			if err := csvWriter.Error(); err != nil {
				log.Printf("Failed to write \"%s\"\n", csvFile.Name())
			}

			saveJson(jsonData, path.Join(jsonPath, simpleId+".json"))
		}
	}

	httpClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	connectUrl := config.Connection.Url + config.Connection.DB

	n := 0
	maxDownloads := config.Connection.DownloadListMaxsize

	for n < maxDownloads {

		req, err := http.NewRequest(http.MethodGet, connectUrl, nil)
		if err != nil {
			log.Panic(err)
		}

		req.Header.Set("Accept", "application/json")

		q := req.URL.Query()
		q.Add("query", config.Connection.Query)
		if len(config.Connection.Fcq) != 0 {
			q.Add("fcq", config.Connection.Fcq)
		}
		q.Add("maximumRecords", strconv.Itoa(config.Connection.DownloadBatchSize))
		q.Add("startRecord", strconv.Itoa(n+1))
		req.URL.RawQuery = q.Encode()

		data, err := downloadJson(&httpClient, req)
		if err != nil {
			log.Panic(err)
		}

		rl := RecordsList{}

		err = json.Unmarshal(data, &rl)
		if err != nil {
			log.Panic(err)
		}

		if rl.NumberOfRecords < maxDownloads {
			maxDownloads = rl.NumberOfRecords
		}

		log.Println(fmt.Sprintf("Start to download [%d-%d]/%d", n+1, n+len(rl.InnerRecordsList.Records), maxDownloads))

		for _, val := range rl.InnerRecordsList.Records {
			n += 1

			downloadUrl := connectUrl + "/" + url.PathEscape(val.RecordIdentifier)

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

			saveToES(jsonData, val.RecordIdentifier)
			saveToFS(jsonData, n, val.RecordIdentifier)
		}

		log.Println(fmt.Sprintf("Downloaded %d/%d. Next record number is %d", n, maxDownloads, rl.NextRecordPosition))
	}
}

func getHref(data []byte) string {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Failed to get href")
		}
	}()

	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Panic(err)
	}

	href := result["children"].([]interface{})[1].(map[string]interface{})["children"].([]interface{})[0].(map[string]interface{})["children"].([]interface{})[3].(map[string]interface{})["children"].([]interface{})[4].(map[string]interface{})["children"].([]interface{})[1].(map[string]interface{})["href"].(string)

	return href
}
