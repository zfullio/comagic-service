package csv

import (
	"fmt"
	csvtag "github.com/artonge/go-csv-tag/v2"
	"log"
	"os"
	"time"
)

func GenerateFile(data any, name string) (err error) {
	csvFile, err := os.Create(name)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvFile.Close()

	options := csvtag.CsvOptions{
		Separator: '|',
	}
	err = csvtag.DumpToFile(data, name, options)
	if err != nil {
		return err
	}
	return err
}

func GenerateFilename(name string) string {
	t := time.Now()
	filename := fmt.Sprintf("%s %s.csv", name, t.Format("2006-01-02 15:04:05"))
	return filename
}
