package cs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type campaignRepository struct {
	db         storage.Client
	BucketName string
	logger     *zerolog.Logger
}

func NewCampaignRepository(client storage.Client, bucketName string, logger *zerolog.Logger) *campaignRepository {
	repoLogger := logger.With().Str("repo", "campaign").Str("type", "cloud-storage").Logger()

	return &campaignRepository{
		db:         client,
		BucketName: bucketName,
		logger:     &repoLogger,
	}
}

func (pr campaignRepository) SendFile(ctx context.Context, filename string) error {
	pr.logger.Trace().Msgf("SendFile: %v", filename)

	bucket := pr.db.Bucket(pr.BucketName)

	err := SendFile(ctx, bucket, filename)
	if err != nil {
		return fmt.Errorf("SendFile: %w", err)
	}

	return nil
}
