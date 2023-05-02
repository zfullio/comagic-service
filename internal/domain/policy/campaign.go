package policy

import (
	"Comagic/internal/domain/service"
	cm "Comagic/pkg/comagic"
	"context"
	"fmt"
)

type CampaignPolicy struct {
	service service.CampaignService
}

func NewCampaignPolicy(service service.CampaignService) *CampaignPolicy {
	return &CampaignPolicy{service: service}
}

func (cp CampaignPolicy) GetCampaigns(ctx context.Context, fields []string, filter cm.Filter, datasetID string, tableID string) error {
	campaigns, err := cp.service.GetCampaigns(fields, filter)
	if err != nil {
		return fmt.Errorf("ошибка получения кампаний: %w", err)
	}

	if len(campaigns) == 0 {
		return fmt.Errorf("кампании | пустой список значений")
	}

	err = cp.service.SendCampaigns(ctx, datasetID, tableID, campaigns)
	if err != nil {
		return fmt.Errorf("ошибка отправки кампаний в bq: %w", err)
	}

	return nil
}

func (cp CampaignPolicy) GetCondition(ctx context.Context, fields []string, filter cm.Filter, datasetID string, tableID string) error {
	campaigns, err := cp.service.GetCampaigns(fields, filter)
	if err != nil {
		return fmt.Errorf("ошибка получения кампаний: %w", err)
	}

	if len(campaigns) == 0 {
		return fmt.Errorf("кампании | пустой список значений")
	}

	err = cp.service.SendCampaignsConditions(ctx, datasetID, tableID, campaigns)
	if err != nil {
		return fmt.Errorf("ошибка отправки campaign conditions в bq: %w", err)
	}

	return nil
}
