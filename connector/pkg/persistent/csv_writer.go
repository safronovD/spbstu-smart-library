package persistent

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

type CSVWriter struct {
	csvWriter *csv.Writer
}

func NewCSVWriter(csvPath string) (*CSVWriter, error) {
	csvFile, err := os.Create(csvPath)
	if err != nil {
		return nil, err
	}
	// TODO how is this file able to be closed?

	csvWriter := csv.NewWriter(csvFile)

	return &CSVWriter{
		csvWriter: csvWriter,
	}, nil
}

func (w *CSVWriter) Write(jsonData []byte, recordId string) {
	if err := w.saveCSV(jsonData, recordId); err != nil {
		log.Printf("Failed to update csv with err: %s", err)
	}
}

func (w *CSVWriter) saveCSV(jsonData []byte, recordId string) error {
	csvLine := []string{recordId, getHref(jsonData)}

	if err := w.csvWriter.Write(csvLine); err != nil {
		return err
	}

	w.csvWriter.Flush()

	if err := w.csvWriter.Error(); err != nil {
		return err
	}

	return nil
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

	href := result["pdfLink"].(string)

	return href
}
