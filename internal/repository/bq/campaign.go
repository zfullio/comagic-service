package bq

import (
	"Comagic/internal/domain/entity"
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type campaignRepository struct {
	db     *bigquery.Client
	table  *bigquery.Table
	logger *zerolog.Logger
}

func NewCampaignRepository(client *bigquery.Client, datasetID string, tableID string, logger *zerolog.Logger) *campaignRepository {
	repoLogger := logger.With().Str("repo", "campaign").Str("type", "bigquery").Logger()

	dataset := client.Dataset(datasetID)
	table := dataset.Table(tableID)

	return &campaignRepository{
		db:     client,
		table:  table,
		logger: &repoLogger,
	}
}

func (cr campaignRepository) SendAny(datasetID string, tableID string, campaigns []entity.Campaign) error {
	campaignsDTO := make([]campaignDTO, 0, len(campaigns))

	for i := 0; i < len(campaigns); i++ {
		item := newCampaignDTO(campaigns[i])
		campaignsDTO = append(campaignsDTO, *item)
	}

	myDataset := cr.db.Dataset(datasetID)
	table := myDataset.Table(tableID)

	u := table.Inserter()
	if err := u.Put(context.Background(), campaignsDTO); err != nil {
		return fmt.Errorf("bigquery error: %w", err)
	}

	return nil
}

func newCampaignDTO(campaign entity.Campaign) *campaignDTO {
	dynamicCallTrackingDTO := DynamicCallTrackingDTO{
		ReservationTime:     campaign.DynamicCallTracking.ReservationTime,
		CountVirtualNumbers: campaign.DynamicCallTracking.CountVirtualNumbers,
		CountVisits:         campaign.DynamicCallTracking.CountVisits,
		CoverageVisitors:    campaign.DynamicCallTracking.CoverageVisitors,
	}

	return &campaignDTO{
		ID:                  campaign.ID,
		Status:              campaign.Status,
		Name:                campaign.Name,
		DynamicCallTracking: dynamicCallTrackingDTO,
		StaticUtmSource:     campaign.StaticUtmSource,
		StaticUtmMedium:     campaign.StaticUtmMedium,
		StaticUtmCampaign:   campaign.StaticUtmCampaign,
		StaticUtmTerm:       campaign.StaticUtmTerm,
		StaticUtmContent:    campaign.StaticUtmContent,
		StaticUtmReferrer:   campaign.StaticUtmReferrer,
		StaticUtmExpid:      campaign.StaticUtmExpid,
		DateUpdate:          campaign.DateUpdate,
	}
}

type campaignDTO struct {
	ID                  int64                  `bigquery:"id"`
	Status              string                 `bigquery:"status"`
	Name                string                 `bigquery:"name"`
	DynamicCallTracking DynamicCallTrackingDTO `bigquery:"dynamic_call_tracking"`
	StaticUtmSource     string                 `bigquery:"static_utm_source"`
	StaticUtmMedium     string                 `bigquery:"static_utm_medium"`
	StaticUtmCampaign   string                 `bigquery:"static_utm_campaign"`
	StaticUtmTerm       string                 `bigquery:"static_utm_term"`
	StaticUtmContent    string                 `bigquery:"static_utm_content"`
	StaticUtmReferrer   string                 `bigquery:"static_utm_referrer"`
	StaticUtmExpid      string                 `bigquery:"static_utm_expid"`
	DateUpdate          time.Time              `bigquery:"date_update"`
}

type DynamicCallTrackingDTO struct {
	ReservationTime     string  `bigquery:"reservation_time"`
	CountVirtualNumbers int64   `bigquery:"count_virtual_numbers"`
	CountVisits         int64   `bigquery:"count_visits"`
	CoverageVisitors    float64 `bigquery:"coverage_visitors"`
}

func (cr campaignRepository) SendCampaignsConditions(datasetID string, tableID string, campaigns []entity.Campaign) error {
	conditions := make([]CampaignConditionDTO, 0)

	for i := 0; i < len(campaigns); i++ {
		for j := 0; j < len(campaigns[i].CampaignConditions); j++ {
			item := newCampaignConditionDTO(campaigns[i].CampaignConditions[j])
			conditions = append(conditions, *item)
		}
	}

	myDataset := cr.db.Dataset(datasetID)
	table := myDataset.Table(tableID)
	u := table.Inserter()

	if err := u.Put(context.Background(), conditions); err != nil {
		return fmt.Errorf("bigquery error: %w", err)
	}

	return nil
}

type CampaignConditionDTO struct {
	ID                int64     `bigquery:"id"`
	GroupID           int       `bigquery:"group_id"`
	Type              string    `bigquery:"type"`
	Value             string    `bigquery:"value"`
	Operator          string    `bigquery:"operator"`
	CampaignParameter string    `bigquery:"campaign_parameter"`
	DateUpdate        time.Time `bigquery:"date_update"`
}

func newCampaignConditionDTO(condition entity.GroupCondition) *CampaignConditionDTO {
	return &CampaignConditionDTO{
		ID:                int64(condition.ID),
		GroupID:           condition.GroupID,
		Type:              condition.Type,
		Value:             condition.Value,
		Operator:          condition.Operator,
		CampaignParameter: condition.CampaignParameter,
		DateUpdate:        condition.DateUpdate,
	}
}
