package persistent

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	json_path = "/tmp/test.json"
	json_dir  = "/tmp/json"
)

type TestStruct struct {
	Field1 string   `json:"field1"`
	Field2 int      `json:"field2"`
	Field3 bool     `json:"field3"`
	Field4 []string `json:"field4"`
}

func TestSaveJSON(t *testing.T) {
	// Create Test struct and set some fields
	testStruct := createTestStruct()

	// Write data to file
	jsonData, err := json.Marshal(testStruct)
	assert.Nil(t, err)

	err = SaveJSON(jsonData, json_path)
	assert.Nil(t, err)

	// Read data from file and check that structures are the same
	readTestStruct, err := readTestStruct(json_path)
	assert.Nil(t, err)
	assert.Equal(t, testStruct, readTestStruct)
}

func TestNewFileSystemWriter(t *testing.T) {
	fsw1 := NewFileSystemWriter(json_dir)
	fsw2 := NewFileSystemWriter(json_dir)

	assert.Equal(t, fsw1, fsw2)
}

func createTestStruct() *TestStruct {
	return &TestStruct{
		Field1: "string",
		Field2: 23,
		Field3: true,
		Field4: []string{"str1", "str2"},
	}
}

func readTestStruct(path string) (*TestStruct, error) {
	readData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	readTestStruct := &TestStruct{}
	if err = json.Unmarshal(readData, readTestStruct); err != nil {
		return nil, err
	}

	return readTestStruct, nil
}
