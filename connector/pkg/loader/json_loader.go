package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/spbstu-smart-library/connector/pkg/config"
	"github.com/spbstu-smart-library/connector/pkg/converter"
	"github.com/spbstu-smart-library/connector/pkg/persistent"
)

type Record struct {
	RecordIdentifier string `json:"recordIdentifier"`
}

type BatchRecord struct {
	RecordsList struct {
		Records []Record `json:"record"`
	} `json:"records"`

	NumberOfRecords    int `json:"numberOfRecords"`
	NextRecordPosition int `json:"nextRecordPosition"`
}

type JsonLoader struct {
	MetaRecordSchema *BatchRecord
	Configuration    *config.JSONConfig
	Writers          []persistent.Writer
	Converter        *converter.JsonConverter
	schema           string
}

func NewJsonLoader(config *config.JSONConfig, outputDirs string) *JsonLoader {
	loader := new(JsonLoader)
	loader.MetaRecordSchema = new(BatchRecord)
	loader.Configuration = config

	if config.Output.ConvertEnable {
		loader.Converter = converter.NewJsonConverter()
		loader.schema = config.Output.ConvertSchema
		if loader.schema == "" {
			loader.Converter = nil
		}
	}

	loader.configureWriters(config, outputDirs)
	return loader
}

func (l *JsonLoader) Download() {
	recordNum := 0
	maxDownloads := l.Configuration.Connection.DownloadListMaxsize
	httpClient := createClient()

	for recordNum < maxDownloads {
		err := l.downloadInfoRecord(httpClient, recordNum)
		if err != nil {
			log.Panic(err)
		}
		if l.MetaRecordSchema.NumberOfRecords < maxDownloads {
			maxDownloads = l.MetaRecordSchema.NumberOfRecords
		}

		log.Println(fmt.Sprintf("Start to download [%d-%d]/%d",
			recordNum+1, recordNum+len(l.MetaRecordSchema.RecordsList.Records), maxDownloads))

		for _, val := range l.MetaRecordSchema.RecordsList.Records {
			recordNum++
			downloadURL := "/" + url.PathEscape(val.RecordIdentifier)

			jsonData, err := l.downloadDataRecord(httpClient, downloadURL)
			if err != nil {
				log.Println(err)
				continue
			}

			if l.Configuration.Output.ConvertEnable && l.Converter != nil {
				convertedJSON, err := l.convertJson(jsonData, l.schema)
				if err != nil {
					log.Printf("JSON converted failed: %s", err)
				}
				jsonData = convertedJSON
			}

			for _, w := range l.Writers {
				w.Write(jsonData, formatID(val.RecordIdentifier))
			}
		}

		log.Printf(fmt.Sprintf("Downloaded %d/%d. Next record number is %d",
			recordNum, maxDownloads, l.MetaRecordSchema.NextRecordPosition))
	}
}

func (l *JsonLoader) downloadInfoRecord(httpClient *http.Client, recordNum int) error {
	request, err := l.createMetaRequest(recordNum)
	if err != nil {
		return err
	}

	data, err := downloadJsonFile(httpClient, request)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, l.MetaRecordSchema)
	if err != nil {
		return err
	}
	return nil
}

func (l *JsonLoader) downloadDataRecord(httpClient *http.Client, JsonURL string) ([]byte, error) {
	request, err := l.createRecordRequest(JsonURL)
	if err != nil {
		return nil, err
	}

	jsonData, err := downloadJsonFile(httpClient, request)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (l *JsonLoader) createMetaRequest(startRecord int) (*http.Request, error) {
	connectURL := l.Configuration.Connection.URL + l.Configuration.Connection.DB
	req, err := http.NewRequest(http.MethodGet, connectURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	query := req.URL.Query()
	query.Add("query", l.Configuration.Connection.Query)

	if len(l.Configuration.Connection.Fcq) != 0 {
		query.Add("fcq", l.Configuration.Connection.Fcq)
	}

	query.Add("maximumRecords", strconv.Itoa(l.Configuration.Connection.DownloadBatchSize))
	query.Add("startRecord", strconv.Itoa(startRecord+1))
	req.URL.RawQuery = query.Encode()
	return req, nil
}

func (l *JsonLoader) createRecordRequest(recordURL string) (*http.Request, error) {
	connectURL := l.Configuration.Connection.URL + l.Configuration.Connection.DB + recordURL
	req, err := http.NewRequest(http.MethodGet, connectURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	query := req.URL.Query()
	query.Add("recordSchema", "gost-7.0.100")
	req.URL.RawQuery = query.Encode()
	return req, nil
}

func (l *JsonLoader) convertJson(jsonData []byte, schema string) ([]byte, error) {
	convertedData, err := l.Converter.Convert(schema, string(jsonData))
	if err != nil {
		return jsonData, err
	}
	return convertedData, nil
}

func (l *JsonLoader) configureWriters(config *config.JSONConfig, outputDir string) {
	if config.Output.Elasticsearch.Enable {
		esConf := config.Output.Elasticsearch
		esWriter, err := persistent.NewESWriter(esConf.Host, esConf.Login, esConf.Login, esConf.Index)
		if err != nil {
			log.Printf("ES connection failed: %s", err)
		} else {
			l.Writers = append(l.Writers, esWriter)
		}
	}

	if config.Output.CsvFile.Enable {
		csvWriter, err := persistent.NewCSVWriter(path.Join(outputDir, config.Output.CsvFile.FileName))
		if err != nil {
			log.Printf("CSV Writer failed: %s", err)
		} else {
			l.Writers = append(l.Writers, csvWriter)
		}
	}

	if config.Output.FileSystem.Enable {
		fsWriter := persistent.NewFileSystemWriter(path.Join(outputDir, config.Output.FileSystem.JSONDir))
		l.Writers = append(l.Writers, fsWriter)
	}
}

func createClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 10 * time.Second
	retryClient.RetryWaitMax = 30 * time.Second
	return retryClient.StandardClient()
}

func downloadJsonFile(client *http.Client, req *http.Request) ([]byte, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != http.StatusMultiStatus) {
		return nil, fmt.Errorf("response failed: %s status code: %d", res.Request.URL, res.StatusCode)
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

func formatID(id string) string {
	return strings.ReplaceAll(id, "\\", "_")
}
