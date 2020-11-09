package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type JsonConfig struct {
	Connection struct {
		Url                 string `yaml:"url"`
		DB                  string `yaml:"db"`
		Query               string `yaml:"query"`
		Fcq                 string `yaml:"fcq"`
		DownloadListMaxsize int    `yaml:"download_list_maxsize"`
		DownloadBatchSize   int    `yaml:"download_batch_size"`
	} `yaml:"connection"`

	Output struct {
		ConvertEnable bool `yaml:"convert_enable"`

		Elasticsearch struct {
			Enable bool   `yaml:"enable"`
			Host   string `yaml:"host"`
			Login  string `yaml:"login"`
			Pwd    string `yaml:"pwd"`
			Index  string `yaml:"index"`
		} `yaml:"elasticsearch"`

		FileSystem struct {
			Enable  bool   `yaml:"enable"`
			JsonDir string `yaml:"json_dir"`
			CsvFile string `yaml:"csv_file"`
		} `yaml:"filesystem"`
	} `yaml:"output"`
}

type PdfConfig struct {
	DB      string `yaml:"db"`
	Dir     string `yaml:"dir"`
	CsvFile string `yaml:"csv_file"`

	Auth struct {
		ASPXAUTH        string `yaml:".ASPXAUTH"`
		ASPNETSessionId string `yaml:"ASP.NET_SessionId"`
	} `yaml:"auth"`
}

type Config struct {
	JsonConf JsonConfig `yaml:"json"`
	PdfConf  PdfConfig  `yaml:"pdf"`
}

func NewConfig(configPath string) (*Config, error) {
	s, err := os.Stat(configPath)
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, fmt.Errorf("'%s' is a directory, not a normal file", configPath)
	}

	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	launchMod := flag.String("launch-mod", "download-json", "a string")
	configFileName := flag.String("config-file", "config.yaml", "a string")
	logFileName := flag.String("log-file", "connector.log", "a string")
	outputDir := flag.String("output-dir", "output", "a string")
	flag.Parse()

	os.Mkdir(*outputDir, os.ModePerm)
	logFilePath := path.Join(".", *outputDir, *logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	config, err := NewConfig(*configFileName)
	if err != nil {
		log.Fatal(err)
	}

	switch *launchMod {
	case "download-json":
		downloadRecords(&config.JsonConf, *outputDir)
	case "download-pdf":
		downloadPdfs(&config.PdfConf, *outputDir)
	case "samples":
		downloadSamples(*outputDir)
	default:
		log.Panic("Launch mod is not correct")
	}
}
