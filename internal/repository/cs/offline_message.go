package cs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type offlineMessageRepository struct {
	db     *storage.Client
	Bucket *storage.BucketHandle
	logger *zerolog.Logger
}

func NewOfflineMessageRepository(client *storage.Client, bucketName string, logger *zerolog.Logger) *offlineMessageRepository {
	repoLogger := logger.With().Str("repo", "offline-message").Str("type", "cloud-storage").Logger()

	bucket := client.Bucket(bucketName)

	return &offlineMessageRepository{
		db:     client,
		Bucket: bucket,
		logger: &repoLogger,
	}
}

func (or offlineMessageRepository) SendFile(ctx context.Context, filename string) error {
	or.logger.Trace().Msgf("SendFile: %v", filename)

	err := SendFile(ctx, or.Bucket, filename)
	if err != nil {
		return fmt.Errorf("SendFile: %w", err)
	}

	return nil
}
