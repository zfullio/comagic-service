package cs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type callRepository struct {
	db         storage.Client
	BucketName string
	logger     *zerolog.Logger
}

func NewCallRepository(client storage.Client, bucketName string, logger *zerolog.Logger) *callRepository {
	repoLogger := logger.With().Str("repo", "call").Str("type", "cloud-storage").Logger()

	return &callRepository{
		db:         client,
		BucketName: bucketName,
		logger:     &repoLogger,
	}
}

func (cr callRepository) SendFile(ctx context.Context, filename string) (err error) {
	cr.logger.Trace().Msgf("SendFile: %v", filename)

	bucket := cr.db.Bucket(cr.BucketName)

	err = SendFile(ctx, bucket, filename)
	if err != nil {
		return fmt.Errorf("SendFile: %w", err)
	}

	return nil
}
