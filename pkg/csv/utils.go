package csv

import (
	"fmt"
	csvtag "github.com/artonge/go-csv-tag/v2"
	"os"
	"time"
)

func GenerateFile(data any, name string) (string, error) {
	t := time.Now()
	pattern := fmt.Sprintf("%s_%s_*.csv", name, t.Format("2006-01-02 15:04:05"))

	csvFile, crErr := os.CreateTemp("", pattern)
	if crErr != nil {
		return "", fmt.Errorf("os.CreateTemp: %w", crErr)
	}

	clErr := csvFile.Close()
	if clErr != nil {
		return "", fmt.Errorf("csvFile.Close: %w", clErr)
	}

	options := csvtag.CsvOptions{
		Separator: '|',
	}

	dErr := csvtag.DumpToFile(data, csvFile.Name(), options)
	if dErr != nil {
		return "filename", fmt.Errorf("csvtag.DumpToFile: %w", dErr)
	}

	return csvFile.Name(), nil
}
