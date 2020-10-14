package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const configFile = "config.yaml"

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

type Config struct {
	Connection struct {
		Url                 string `yaml:"url"`
		DB                  string `yaml:"db"`
		Query               string `yaml:"query"`
		Fcq                 string `yaml:"fcq"`
		DownloadListMaxsize int    `yaml:"download_list_maxsize"`
		DownloadBatchSize   int    `yaml:"download_batch_size"`
	} `yaml:"connection"`

	Output struct {
		OutputDir  string `yaml:"common_dir"`
		SamplesDir string `yaml:"samples_dir"`
	} `yaml:"output"`

	Auth struct {
		ASPXAUTH        string `yaml:".ASPXAUTH"`
		ASPNETSessionId string `yaml:"ASP.NET_SessionId"`
	}
}

func NewConfig(configPath string) (*Config, error) {
	s, err := os.Stat(configPath)
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, fmt.Errorf("'%s' is a directory, not a normal file", configPath)
	}

	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

/*func authorize(client *http.Client, name string, password string) http.Cookie {

	url :="https://cas.spbstu.ru/login"

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("username", name)
	q.Add("password", password)
	q.Add("execution", "07385315-351d-4938-ab3c-f556d3880dad_ZXlKaGJHY2lPaUpJVXpVeE1pSjkuMXB3czNaVk1xbTlITS80QUd6Y2phVytXYzlGRGdQQWRQT2NVaTVkb29kd0tOL2FKaVkwWUNUVHV4TXRLRmhVK3lkSFdTeUx1Nm44VnhMSTk0TlYzK1ZGK3F1ZkZnMkFjNDJVYTZIM21XTDlDVmRtM1hkQnJYT3lYWDBUZ1cvT1R0V2c4akJRTTVnUGVvQjhldjZPZjFuMDR5OElYZXhWdFRoUDJNVVZSWitMZHhoQjM4REpOZERKd0NES1V5OG9PSHJYYjJEQm5IbHEwYlhFdUd4UVlYN21lS29LQWFHdlJSWHExeW5ydTQ2aGF5M0hBUEhuZ0toL0dPZ3VTUjZ6RUtxajZQaE5UUExyclFoV3IrdUN5UlFyc0NsZGM4V2lnWDJPYm51Rnk2N0xMblVtbmF2WDBvYmNoakFqWVd3eVFkTm8wNXJ4d0JVVDQ3Mitzd1FMWEdmU3NyZFdtMDhtZXVqV09EYWRBWmZOZGtqMVQzZE9jL1lGOS81cFpDSHJHRFczd29ZSXlGaWtkSUdyZ0RtK0Q1bWREaWVGWVR0a21YQXVmNU81L0lLaHFFUmhMenZVYkZtU0J3RjBpVFpvK3g0NlpIcDN3QXdrVUdvQ3pocm56L1d6dGh5clh2M0llRHVGTDdkU3d3Y2QvblZIZFBmclpXaDRWZW9qY0NENlhpQkhENzNkTi91NUFGdTBCR1o2WUx4WElIdUduNGtvVE5MbENzc3Qxd0NobVM2ZkFEbWs0WEFkeGVUMStmSjA2YmNlTGUvcTVtb1ljaFhncXA2Y2IybjhIZUY0RHNrUmFaOGtMdEdiWm5jNHpoZ2pxR1FKRWFOVjRheGdjcmNpR1AwMyttVzlCempteHV5ek54eDNwVDhKYmI1Q1A1WEhGS25WZ0FFYU1kamhLZU9WZlF6R0ZyUGNSK1ozemRrczJYWXg2blhaZVdaWkdDSzVSZHc4U2w2dTJ2R21MZy9wZ21XR0hPKzczTkVlVVkyVGNTYnl1eW05RHNYQ0twcHBKaGFDclJpUEREbmMxcFdRMFFpaVZQU01uNHRFWlRlb2VhS3RDNjVoczZzTUttL3huTEU3NHRQSzNod0srYk5MRElQbWwwaGZtVDNVZlVTYzBvL1pxZ20xZUVTMUFvZzYxRG9FbzZvLzhwMjBnL0tHVHVpb3pHZmVickYwTVJuTnpuMDRQV2NnRFRmcU9mazVuOWIweWV1eFNmeUJuYUxwZkVMbVpSOXpSOXRzWVJBaWR3UmNaRk1pSFNJQkYzZ2lTUzFLMkczZDQrOGlsMkt4TUpUVkc5WU9oY3pzUGlXa3p2cmg2RTRYcWtJOD0uaXhXeU9wSDE1YVVkYW0xM194TlE1TXk1WVBOaF9KR05rVzJWaUJzM3NBRzNzdE1BMHZuOXZDbnpMRzJhNC0wRU1NdEhGcWFRRXZ0cnE1eDNuQTBtN2c=")
	q.Add("_eventId", "submit")
	q.Add("geolocation", "")

	req.URL.RawQuery = q.Encode()

	rsp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp)

	return *rsp.Cookies()[0]
}*/

func getHref(data []byte) string {

	defer func() {
		if recover() != nil {
			log.Println("Failed to get href")
		}
	}()

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	href := result["children"].([]interface{})[1].(map[string]interface{})["children"].([]interface{})[0].(map[string]interface{})["children"].([]interface{})[3].(map[string]interface{})["children"].([]interface{})[4].(map[string]interface{})["children"].([]interface{})[1].(map[string]interface{})["href"].(string)

	return href
}

func downloadRecords(config *Config, client *http.Client) {
	connectUrl := config.Connection.Url + config.Connection.DB
	aspxauth := config.Auth.ASPXAUTH
	sessionId := config.Auth.ASPNETSessionId

	var cookies [2]http.Cookie
	cookies[0] = http.Cookie{Name: ".ASPXAUTH", Value: aspxauth}
	cookies[1] = http.Cookie{Name: "ASP.NET_SessionId", Value: sessionId}

	commonPath := path.Join(".", config.Output.OutputDir, config.Connection.DB)
	jsonPath := path.Join(commonPath, "jsons")
	pdfPath := path.Join(commonPath, "pdfs")
	logPath := path.Join(commonPath, "log")

	os.Mkdir(commonPath, os.ModePerm)
	os.Mkdir(jsonPath, os.ModePerm)
	os.Mkdir(pdfPath, os.ModePerm)
	os.Mkdir(logPath, os.ModePerm)

	f, err := os.Create(path.Join(commonPath, "id_list.csv"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n := 0
	maxDownloads := config.Connection.DownloadListMaxsize

	for n < maxDownloads {

		req, err := http.NewRequest(http.MethodGet, connectUrl, nil)
		if err != nil {
			panic(err)
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

		data, err := downloadJson(client, req, path.Join(logPath, "records_list-log"+strconv.Itoa(n)+".json"))
		if err != nil {
			log.Fatal(err)
		}

		rl := RecordsList{}

		err = json.Unmarshal(data, &rl)
		if err != nil {
			log.Fatal(err)
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
				log.Fatal(err)
			}

			req.Header.Set("Accept", "application/json")

			q := req.URL.Query()
			q.Add("recordSchema", "gost-7.0.100")
			req.URL.RawQuery = q.Encode()

			jsonData, err := downloadJson(client, req, path.Join(jsonPath, strconv.Itoa(n)+".json"))
			if err != nil {
				log.Println(err)
				continue
			}

			//TODO: use "encoding/csv"
			f.Write([]byte(strconv.Itoa(n)))
			f.Write([]byte(","))
			f.Write([]byte(val.RecordIdentifier))

			href := getHref(jsonData)

			if href != "" {
				splitedHref := strings.Split(href, "/")
				pdfName := splitedHref[len(splitedHref)-1]

				f.Write([]byte(","))
				f.Write([]byte(pdfName))

				downloadFile(client, cookies, href+"/download", path.Join(pdfPath, pdfName))
			} else {
				log.Println("Failed to get href")
			}

			f.Write([]byte("\n"))
		}

		log.Println(fmt.Sprintf("Downloaded %d/%d. Next record number is %d", n, maxDownloads, rl.NextRecordPosition))
	}
}

func downloadSamples(config *Config, client *http.Client) {
	connectUrl := config.Connection.Url

	commonPath := path.Join(".", config.Output.OutputDir, config.Output.SamplesDir)
	os.Mkdir(config.Output.OutputDir, os.ModePerm)
	os.Mkdir(commonPath, os.ModePerm)

	req, err := http.NewRequest("PROPFIND", connectUrl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")

	data, err := downloadJson(client, req, path.Join(commonPath, "db_list"+".json"))
	if err != nil {
		log.Fatal(err)
	}

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

		_, err = downloadJson(client, req, path.Join(dbPath, "db_index"+".json"))
		if err != nil {
			log.Fatal(err)
		}

		//TODO: without config editing
		config.Output.OutputDir = commonPath
		config.Connection.DB = dbName
		config.Connection.Query = "(bib.volume=*)"
		config.Connection.DownloadListMaxsize = 2
		config.Connection.DownloadBatchSize = 2

		downloadRecords(config, client)
	}
}

func main() {
	logFile, err := os.OpenFile("connector.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	//downloadRecords(config, &spaceClient)
	downloadSamples(config, &spaceClient)
}
