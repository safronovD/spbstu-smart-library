package loader

import (
	"encoding/csv"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/spbstu-smart-library/connector/pkg/config"
)

type PDFLoader struct {
	csvPath string
	pdfDir  string
	cookies [2]http.Cookie

	client http.Client
}

func NewPDFLoader(pdfConfig *config.PDFConfig, outputDir string) *PDFLoader {
	csvpath, pdfdir := configureDirs(pdfConfig, outputDir)
	cookies := createCookies(pdfConfig)
	client := http.Client{Timeout: time.Second * 60}
	return &PDFLoader{csvPath: csvpath, pdfDir: pdfdir, cookies: cookies, client: client}
}

func (l *PDFLoader) Download() {
	csvfile, err := os.Open(l.csvPath)
	if err != nil {
		log.Panic(err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	recordList, err := reader.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	for _, record := range recordList {
		err := l.downloadRecord(record[1], path.Join(l.pdfDir, record[0]+".pdf"))
		if err != nil {
			log.Print(err)
		}
	}
}

func (l *PDFLoader) downloadRecord(url string, fileName string) error {
	request, err := l.createRequest(url)
	if err != nil {
		return err
	}

	response, err := l.client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("Failed to GET " + url + " - Status code: " + response.Status)
	}

	err = writeResponse(response, fileName)
	if err != nil {
		return err
	}
	log.Printf("%v dowloaded", url)
	return nil
}

func (l *PDFLoader) createRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range l.cookies {
		request.AddCookie(&cookie)
	}
	//TODO: Check ACCEPT header for minimize ready-to-get formats
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Encoding", "gzip, deflate, bt")
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("Connection", "keep-alive")
	return request, nil
}

func configureDirs(pdfConfig *config.PDFConfig, outputDir string) (string, string) {
	pdfPath := path.Join(".", outputDir, pdfConfig.DB, pdfConfig.Dir)
	if err := os.Mkdir(pdfPath, os.ModePerm); err != nil {
		log.Panic(err)
	}
	csvPath := path.Join(".", outputDir, pdfConfig.CsvFile)

	return pdfPath, csvPath
}

func createCookies(pdfConfig *config.PDFConfig) [2]http.Cookie {
	aspxauth := http.Cookie{Name: ".ASPXAUTH", Value: pdfConfig.Auth.ASPXAUTH}
	sessionID := http.Cookie{Name: "ASP.NET_SessionId", Value: pdfConfig.Auth.ASPNETSessionID}

	return [2]http.Cookie{aspxauth, sessionID}
}

func writeResponse(response *http.Response, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		os.Remove(fileName)
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
