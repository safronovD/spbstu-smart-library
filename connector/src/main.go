package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path"

	"github.com/spbstu-smart-library/connector/pkg"
)

func main() {
	launchMod := flag.String("launch-mod", "download-pdf", "a string")
	configFileName := flag.String("config-file", "config.yaml", "a string")
	logFileName := flag.String("log-file", "connector.log", "a string")
	outputDir := flag.String("output-dir", "output", "a string")
	flag.Parse()

	if err := os.Mkdir(*outputDir, os.ModePerm); err != nil {
		log.Panic(err)
	}

	logFilePath := path.Join(".", *outputDir, *logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	config, err := pkg.NewConfig(*configFileName)
	if err != nil {
		log.Panic(err)
	}

	switch *launchMod {
	case "download-json":
		pkg.DownloadRecords(&config.JSONConf, *outputDir)
	case "download-pdf":
		pkg.DownloadPDFFiles(&config.PDFConf, *outputDir)
	case "samples":
		pkg.DownloadSamples(*outputDir)
	default:
		log.Panic("Launch mod is not correct")
	}
}
