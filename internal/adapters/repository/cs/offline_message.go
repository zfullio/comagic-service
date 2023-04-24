package cs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type offlineMessageRepository struct {
	db         storage.Client
	BucketName string
	logger     *zerolog.Logger
}

func NewOfflineMessageRepository(client storage.Client, bucketName string, logger *zerolog.Logger) *offlineMessageRepository {
	repoLogger := logger.With().Str("repo", "offline-message").Str("type", "cloud-storage").Logger()

	return &offlineMessageRepository{
		db:         client,
		BucketName: bucketName,
		logger:     &repoLogger,
	}
}

func (or offlineMessageRepository) SendFile(ctx context.Context, filename string) (err error) {
	or.logger.Trace().Msgf("SendFile: %v", filename)

	bucket := or.db.Bucket(or.BucketName)

	err = SendFile(ctx, bucket, filename)
	if err != nil {
		return fmt.Errorf("SendFile: %w", err)
	}

	return nil
}
