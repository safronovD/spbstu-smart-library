package converter

import (
	"encoding/json"
)

type JsonConverter struct {
	filler *jsonFiller
	parser *jsonParser
}

func NewJsonConverter() *JsonConverter {
	return &JsonConverter{
		filler: NewJsonFiller(),
		parser: NewJsonParser(),
	}
}

func (c *JsonConverter) Convert(schema string, rawData string) ([]byte, error) {
	c.parser.Parse(rawData)
	filledJson := c.filler.Fill(schema, c.parser)
	formattedData, err := json.MarshalIndent(filledJson, " ", " ")
	if err != nil {
		return nil, err
	}
	return formattedData, nil
}
