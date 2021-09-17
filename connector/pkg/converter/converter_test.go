package converter

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

const (
	jsonFile       = ".\\test_files\\test_record.json"
	newJsonFile    = ".\\test_files\\result_record.json"
	prettyJsonFile = ".\\test_files\\result_for_test_record.json"
	schema         = "author, keyWords, links, worktype"
)

func TestJsonConverter(t *testing.T) {
	converter := NewJsonConverter()
	err := converter.Convert(schema, jsonFile, newJsonFile)
	if err != nil {
		t.Error(err)
	}

	file1, err := os.Open(newJsonFile)
	if err != nil {
		t.Error(err)
	}
	defer file1.Close()
	data1, err := ioutil.ReadAll(file1)

	file2, err := os.Open(prettyJsonFile)
	if err != nil {
		t.Error(err)
	}
	defer file2.Close()
	data2, err := ioutil.ReadAll(file2)

	assert.ElementsMatch(t, data1, data2)
}
