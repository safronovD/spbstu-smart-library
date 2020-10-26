package main

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func downloadFile(client *http.Client, cookies [2]http.Cookie, url string, path string) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Failed to download " + url)
		}
	}()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic(err)
	}

	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, bt")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	rsp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if rsp.StatusCode != http.StatusOK {
		log.Panic(errors.New("Failed to GET " + url + " Code: " + rsp.Status))
	}

	file, err := os.Create(path)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		os.Remove(path)
		log.Panic(err)
	}

	file.Write(data)

	log.Println(url + " downloaded")
}

func downloadPdfs(config *PdfConfig, outputDir string) {
	commonPath := path.Join(".", outputDir, config.DB)
	pdfPath := path.Join(commonPath, config.Dir)

	if _, err := os.Stat(commonPath); err != nil {
		log.Panic(err)
	}

	os.Mkdir(pdfPath, os.ModePerm)

	//TODO: Set in config full path to csv
	csvFile, err := os.Open(path.Join(commonPath, config.CsvFile))
	if err != nil {
		log.Panic(err)
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)

	httpClient := http.Client{
		Timeout: time.Second * 60, // Timeout after 2 seconds
	}

	aspxauth := config.Auth.ASPXAUTH
	sessionId := config.Auth.ASPNETSessionId

	var cookies [2]http.Cookie
	cookies[0] = http.Cookie{Name: ".ASPXAUTH", Value: aspxauth}
	cookies[1] = http.Cookie{Name: "ASP.NET_SessionId", Value: sessionId}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}

		downloadFile(&httpClient, cookies, record[3], path.Join(pdfPath, record[2]+".pdf"))
	}
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
