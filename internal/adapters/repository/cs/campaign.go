package cs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"time"
)

type campaignRepository struct {
	db         storage.Client
	BucketName string
}

func NewCampaignRepository(client storage.Client, bucketName string) *campaignRepository {
	return &campaignRepository{
		db:         client,
		BucketName: bucketName,
	}
}

func (pr campaignRepository) SendFile(ctx context.Context, filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	o := pr.db.Bucket(pr.BucketName).Object(filename)
	o = o.If(storage.Conditions{DoesNotExist: true})
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	fmt.Printf("%v uploaded.\n", filename)
	if err != nil {
		return errors.Errorf("statFile: unable to stat file from bucket %q, file %q: %v", pr.BucketName, filename, err)
	}
	return nil
}
