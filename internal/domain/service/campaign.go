package service

import (
	"Comagic/internal/domain/entity"
	"Comagic/pkg/comagic"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type CampaignRepositoryTracking interface {
	GetWithFilter(fields []string, filter comagic.Filter) (campaigns []entity.Campaign, err error)
}

type CampaignRepositoryBQ interface {
	SendAny(ctx context.Context, datasetID string, tableID string, campaigns []entity.Campaign) (err error)
	SendCampaignsConditions(ctx context.Context, datasetID string, tableID string, campaigns []entity.Campaign) (err error)
}

type CampaignService struct {
	tracking CampaignRepositoryTracking
	bq       CampaignRepositoryBQ
	logger   *zerolog.Logger
}

func NewCampaignService(tracking CampaignRepositoryTracking, bq CampaignRepositoryBQ, logger *zerolog.Logger) *CampaignService {
	serviceLogger := logger.With().Str("service", "campaign").Logger()

	return &CampaignService{
		tracking: tracking,
		bq:       bq,
		logger:   &serviceLogger,
	}
}

func (s *CampaignService) GetCampaigns(fields []string, filter comagic.Filter) (campaigns []entity.Campaign, err error) {
	s.logger.Info().Msg("GetCampaigns")

	campaigns, err = s.tracking.GetWithFilter(fields, filter)
	if err != nil {
		return campaigns, fmt.Errorf("ошибка получения кампаний: %w", err)
	}

	return campaigns, nil
}

func (s *CampaignService) SendCampaigns(ctx context.Context, datasetID string, tableID string, campaigns []entity.Campaign) error {
	s.logger.Info().Msg("SendCampaigns")

	err := s.bq.SendAny(ctx, datasetID, tableID, campaigns)
	if err != nil {
		return fmt.Errorf("ошибка отправки данных %w", err)
	}

	return nil
}

func (s *CampaignService) SendCampaignsConditions(ctx context.Context, datasetID string, tableID string, campaigns []entity.Campaign) error {
	s.logger.Info().Msg("SendCampaignsConditions")

	err := s.bq.SendCampaignsConditions(ctx, datasetID, tableID, campaigns)
	if err != nil {
		return fmt.Errorf("ошибка отправки данных %w", err)
	}

	return nil
}
