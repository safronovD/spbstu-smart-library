package converter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type jsonConverter struct {
	filler *jsonFiller
	parser *jsonParser
}

func NewJsonConverter() *jsonConverter {
	return &jsonConverter{
		filler: NewJsonFiller(),
		parser: NewJsonParser(),
	}
}

func (c *jsonConverter) Convert(schema string, jsonPath string, newJsonPath string) error {
	data, err := readJson(jsonPath)
	if err != nil {
		return err
	}
	c.parser.Parse(string(data))
	fmt.Println(c.parser.fields)
	filledJson := c.filler.Fill(schema, c.parser)
	fmt.Println(c.filler.newJson)

	bytes, err := json.MarshalIndent(filledJson, " ", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	err = writeJson(newJsonPath, bytes)
	if err != nil {
		return err
	}
	return nil
}

func readJson(jsonPath string) ([]byte, error) {
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func writeJson(jsonPath string, data []byte) error {
	file, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
