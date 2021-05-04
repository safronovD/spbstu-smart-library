package persistent

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

const (
	json_path = "/tmp/test.json"
)

type TestStruct struct {
	Field1 string   `json:"field1"`
	Field2 int      `json:"field2"`
	Field3 bool     `json:"field3"`
	Field4 []string `json:"field4"`
}

func TestSaveJSON(t *testing.T) {
	// Create Config struct and set some fields
	testStruct := &TestStruct{}

	testStruct.Field1 = "string"
	testStruct.Field2 = 23
	testStruct.Field3 = true
	testStruct.Field4 = []string{"str1", "str2"}

	// Write data to file
	jsonData, err := json.Marshal(testStruct)
	assert.Nil(t, err)

	err = SaveJSON(jsonData, json_path)
	assert.Nil(t, err)

	// Read data from file with NewConfig method and check that structures are the same
	readData, err := ioutil.ReadFile(json_path)
	assert.Nil(t, err)

	readTestStruct := &TestStruct{}
	err = json.Unmarshal(readData, readTestStruct)
	assert.Nil(t, err)
	assert.Equal(t, testStruct, readTestStruct)
}
