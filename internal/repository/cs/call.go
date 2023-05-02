package cs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type callRepository struct {
	db     *storage.Client
	Bucket *storage.BucketHandle
	logger *zerolog.Logger
}

func NewCallRepository(client *storage.Client, bucketName string, logger *zerolog.Logger) *callRepository {
	repoLogger := logger.With().Str("repo", "call").Str("type", "cloud-storage").Logger()

	bucket := client.Bucket(bucketName)

	return &callRepository{
		db:     client,
		Bucket: bucket,
		logger: &repoLogger,
	}
}

func (cr callRepository) SendFile(ctx context.Context, filename string) error {
	cr.logger.Trace().Msgf("SendFile: %v", filename)

	err := SendFile(ctx, cr.Bucket, filename)
	if err != nil {
		return fmt.Errorf("SendFile: %w", err)
	}

	return nil
}
