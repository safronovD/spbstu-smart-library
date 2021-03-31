package loader

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/hashicorp/go-retryablehttp"
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

func saveJSON(data []byte, path string) {
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
	if err != nil {
		log.Panic(err)
	}

	_, err = prettyJSON.WriteTo(jsonFile)

	if err != nil {
		os.Remove(jsonFile.Name())
		log.Panic(err)
	} else {
		log.Println(fmt.Sprintf("Json file \"%s\" saved", jsonFile.Name()))
	}
}

func downloadJSON(client *http.Client, req *http.Request) ([]byte, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != http.StatusMultiStatus) {
		return nil, fmt.Errorf("Response failed: %s Status code: %d", res.Request.URL, res.StatusCode)
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

func DownloadRecords(config *JSONConfig, outputDir string) {
	saveToES := func(jsonData []byte, recordId string) {}
	saveToFS := func(jsonData []byte, recordId string) {}
	saveCSV := func(jsonData []byte, recordId string) {}
	convertJSONData := func(jsonData *[]byte) {}

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

		ctx := context.Background()

		saveToES = func(jsonData []byte, recordId string) {
			defer func() {
				if err := recover(); err != nil {
					log.Println("Failed to load data to ES")
				}
			}()

			req := esapi.IndexRequest{
				Index:      config.Output.Elasticsearch.Index,
				DocumentID: recordId,
				Body:       bytes.NewReader(jsonData),
				Refresh:    "true",
			}

			res, err := req.Do(ctx, es)
			if err != nil || res.StatusCode >= 300 || res.StatusCode < 200 {
				log.Panicf("IndexRequest ERROR: %s, %s", err, res)
			}
			defer res.Body.Close()

			log.Printf("Record with id \"%s\" send to ES.\n", recordId)
		}
	}

	if config.Output.FileSystem.Enable {
		commonPath := path.Join(".", outputDir, config.Connection.DB)
		jsonPath := path.Join(commonPath, config.Output.FileSystem.JSONDir)

		if _, err := os.Stat(commonPath); os.IsNotExist(err) {
			if err := os.Mkdir(commonPath, os.ModePerm); err != nil {
				log.Panic(err)
			}
		}

			if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		if err := os.Mkdir(jsonPath, os.ModePerm); err != nil {
			log.Panic(err)
		}

		saveToFS = func(jsonData []byte, recordId string) {
			saveJSON(jsonData, path.Join(jsonPath, recordId+".json"))
		}
	}

	csvFile, err := os.Create(path.Join(outputDir, config.Output.FileSystem.CsvFile))
	if err != nil {
		log.Panic(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	saveCSV = func(jsonData []byte, recordId string) {
		csvLine := []string{recordId, getHref(jsonData)}

		if err := csvWriter.Write(csvLine); err != nil {
			log.Printf("Failed to write \"%s\"\n", csvFile.Name())
		}

		csvWriter.Flush()

		if err := csvWriter.Error(); err != nil {
			log.Printf("Failed to write \"%s\"\n", csvFile.Name())
		}
	}

	if config.Output.ConvertEnable {
		convertJSONData = func(jsonData *[]byte) {
			// Poor logic, saving temporart to file to convert json. Change it later
			saveJSON(*jsonData, "./tmp1.json")

			cmd := exec.Command("python3", "./utils/json_converter3.py", "./tmp1.json", "./res.json")

			if err = cmd.Run(); err != nil {
				log.Println(err)
			}

			convertedJSON, err := ioutil.ReadFile("./res.json")
			if err != nil {
				log.Println(err)
				log.Println("Json convert phase failed")
			} else {
				*jsonData = convertedJSON
			}
		}
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 10 * time.Second
	retryClient.RetryWaitMax = 30 * time.Second
	httpClient := retryClient.StandardClient()

	connectURL := config.Connection.URL + config.Connection.DB

	n := 0
	maxDownloads := config.Connection.DownloadListMaxsize

	for n < maxDownloads {
		req, err := http.NewRequest(http.MethodGet, connectURL, nil)
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

		data, err := downloadJSON(httpClient, req)
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
			n++

			downloadURL := connectURL + "/" + url.PathEscape(val.RecordIdentifier)

			req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
			if err != nil {
				log.Panic(err)
			}

			req.Header.Set("Accept", "application/json")

			q := req.URL.Query()
			q.Add("recordSchema", "gost-7.0.100")
			req.URL.RawQuery = q.Encode()

			jsonData, err := downloadJSON(httpClient, req)
			if err != nil {
				log.Println(err)
				continue
			}

			formattedID := strings.ReplaceAll(val.RecordIdentifier, "\\", "_")

			convertJSONData(&jsonData)

			saveToES(jsonData, formattedID)
			saveToFS(jsonData, formattedID)
			saveCSV(jsonData, formattedID)
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

	if err := json.Unmarshal(data, &result); err != nil {
		log.Panic(err)
	}

	href := result["pdfLink"].(string)

	return href
}
