package cs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func SendFile(ctx context.Context, bucket *storage.BucketHandle, filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	items := strings.Split(filename, "/")
	o := bucket.Object(items[len(items)-1])
	o = o.If(storage.Conditions{DoesNotExist: true})
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	if err != nil {
		return err
	}

	return nil
}
