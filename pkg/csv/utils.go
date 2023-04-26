package csv

import (
	"fmt"
	csvtag "github.com/artonge/go-csv-tag/v2"
	"log"
	"os"
	"time"
)

func GenerateFile(data any, name string) (filename string, err error) {
	t := time.Now()
	pattern := fmt.Sprintf("%s_%s_*.csv", name, t.Format("2006-01-02 15:04:05"))

	csvFile, err := os.CreateTemp("", pattern)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	err = csvFile.Close()
	if err != nil {
		return "", err
	}

	options := csvtag.CsvOptions{
		Separator: '|',
	}
	err = csvtag.DumpToFile(data, csvFile.Name(), options)
	if err != nil {
		return "", err
	}

	return csvFile.Name(), nil
}

func ParseFile(data any, name string) (filename string, err error) {
	t := time.Now()
	pattern := fmt.Sprintf("%s_%s_*.csv", name, t.Format("2006-01-02 15:04:05"))

	csvFile, err := os.CreateTemp("", pattern)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	err = csvFile.Close()
	if err != nil {
		return "", err
	}

	options := csvtag.CsvOptions{
		Separator: '|',
	}
	err = csvtag.DumpToFile(data, csvFile.Name(), options)
	if err != nil {
		return "", err
	}

	return csvFile.Name(), nil
}

type Demo struct { // A structure with tags
	Client           string `csv:"client"`
	ComagicToken     string `csv:"comagic_token"`
	GoogleServiceKey string `csv:"google_service_key"`
	ProjectId        string `csv:"project_id"`
	DatasetId        string `csv:"dataset_id"`
	BucketName       string `csv:"bucket_name"`
}
