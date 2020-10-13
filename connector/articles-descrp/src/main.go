package main

import (
	"bufio"
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

func downloadFile(client *http.Client, cookies [2]http.Cookie, url string, path string) {

	defer func() {
		if recover() != nil {
			fmt.Println("Failed to download " + url)
			os.Remove(path)
		}
	}()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
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
		fmt.Println("Failed to download " + url + " Code: " + rsp.Status)
		return
	}

	file, err := os.Create(path)
	if err != nil {
		panic(1)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(1)
	}

	file.Write(data)

	fmt.Println(url + " downloaded")
}

func getHref(data []byte) string {

	defer func() {
		if recover() != nil {
			fmt.Println("Failed to get href")
		}
	}()

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	href := result["children"].([]interface{})[1].(map[string]interface{})["children"].([]interface{})[0].(map[string]interface{})["children"].([]interface{})[3].(map[string]interface{})["children"].([]interface{})[4].(map[string]interface{})["children"].([]interface{})[1].(map[string]interface{})["href"].(string)

	return href
}

func main() {

	db := "EBOOKS"
	url := "https://ruslan.library.spbstu.ru/rrs-web/db/" + db

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Login to https://elib.spbstu.ru/ and see your .ASPXAUTH and ASP.NET_SessionId cookies" +
		"(Pure implementation, we will develop normal authorization later :) )")
	fmt.Print(".ASPXAUTH: ")
	aspxauth, _ := reader.ReadString('\n')
	fmt.Print("ASP.NET_SessionId: ")
	sessionId, _ := reader.ReadString('\n')

	aspxauth = aspxauth[0:len(aspxauth) - 1]
	sessionId = sessionId[0:len(sessionId) - 1]

	var cookies [2]http.Cookie
	cookies[0] = http.Cookie{Name: ".ASPXAUTH", Value: aspxauth}
	cookies[1] = http.Cookie{Name: "ASP.NET_SessionId", Value: sessionId}

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	os.Mkdir("./dir", os.ModePerm)
	os.Mkdir("./dir/jsons", os.ModePerm)
	os.Mkdir("./dir/pdfs", os.ModePerm)

	f, err := os.Create("./dir/ids.txt")
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

			jsonFile, err := os.Create("./dir/jsons/" + strconv.Itoa(n) + ".json")
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

			href := getHref(data)

			if href != "" {
				splitedHref := strings.Split(href, "/")
				pdfName := splitedHref[len(splitedHref) - 1]

				downloadFile(&spaceClient, cookies, href + "/download", "./dir/pdfs/" + pdfName)
			} else {
				fmt.Println("Failed to get href")
			}
			n += 1
		}
	}
}
