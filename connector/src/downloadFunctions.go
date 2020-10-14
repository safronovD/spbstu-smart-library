package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func downloadFile(client *http.Client, cookies [2]http.Cookie, url string, path string) {

	defer func() {
		if recover() != nil {
			log.Println("Failed to download " + url)
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
		log.Println("Failed to download " + url + " Code: " + rsp.Status)
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

	log.Println(url + " downloaded")
}

func downloadJson(client *http.Client, req *http.Request, path string) ([]byte, error) {

	jsonFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode != 200) && (res.StatusCode != 207) {
		return nil, errors.New(fmt.Sprintf("Response status code is %d", res.StatusCode))
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

	//TODO: without double memory
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, data, "", "\t")
	_, err = jsonFile.Write(prettyJSON.Bytes())
	if err != nil {
		os.Remove(jsonFile.Name())
		return nil, err
	} else {
		log.Println(fmt.Sprintf("json file \"%s\" writed", jsonFile.Name()))
	}

	return data, nil
}
