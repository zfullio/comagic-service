package comagic

import (
	"Comagic/internal/domain/entity"
	"Comagic/pkg/comagic"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type campaignRepository struct {
	client comagic.Client
	logger *zerolog.Logger
}

func NewCampaignRepository(tracking comagic.Client, logger *zerolog.Logger) *campaignRepository {
	cmLogger := logger.With().Str("repo", "campaign").Str("type", "comagic").Logger()

	return &campaignRepository{
		client: tracking,
		logger: &cmLogger,
	}
}

func (cr *campaignRepository) GetWithFilter(fields []string, filter comagic.Filter) (campaigns []entity.Campaign, err error) {
	cr.logger.Trace().Msg("GetWithFilter")

	data, err := cr.client.GetCampaigns(fields, filter)
	if err != nil {
		return campaigns, fmt.Errorf("ошибка получения кампаний: %w", err)
	}

	campaignsRepo := data.Result.Data
	t := time.Now()

	for i := 0; i < len(campaignsRepo); i++ {
		item := newCampaign(campaignsRepo[i], t)
		campaigns = append(campaigns, *item)
	}

	return campaigns, err
}

func newCampaign(campaign comagic.CampaignInfo, dateUpdate time.Time) *entity.Campaign {
	dynamicCallTracking := entity.DynamicCallTracking{
		ReservationTime:     campaign.DynamicCallTracking.ReservationTime,
		CountVirtualNumbers: campaign.DynamicCallTracking.CountVirtualNumbers,
		CountVisits:         campaign.DynamicCallTracking.CountVisits,
		CoverageVisitors:    campaign.DynamicCallTracking.CoverageVisitors,
	}
	conditions := make([]entity.GroupCondition, 0)

	for idx, group := range campaign.CampaignConditions.GroupConditions {
		for _, item := range group {
			condition := entity.GroupCondition{
				ID:                int(campaign.ID),
				GroupID:           idx,
				Type:              item.Type,
				Value:             item.Value,
				Operator:          item.Operator,
				CampaignParameter: item.CampaignParameter,
				DateUpdate:        dateUpdate,
			}
			conditions = append(conditions, condition)
		}
	}

	return &entity.Campaign{
		ID:                  campaign.ID,
		Status:              campaign.Status,
		CreationTime:        campaign.CreationTime,
		Description:         campaign.Description,
		SiteId:              campaign.SiteId,
		SiteDomainName:      campaign.SiteDomainName,
		Costs:               campaign.Costs,
		CostRatio:           campaign.CostRatio,
		CostRatioOperator:   campaign.CostRatioOperator,
		Engine:              campaign.Engine,
		Type:                campaign.Type,
		Name:                campaign.Name,
		StaticUtmSource:     campaign.StaticUtmSource,
		StaticUtmMedium:     campaign.StaticUtmMedium,
		StaticUtmCampaign:   campaign.StaticUtmCampaign,
		StaticUtmTerm:       campaign.StaticUtmTerm,
		StaticUtmContent:    campaign.StaticUtmContent,
		StaticUtmReferrer:   campaign.StaticUtmReferrer,
		StaticUtmExpid:      campaign.StaticUtmExpid,
		DynamicCallTracking: dynamicCallTracking,
		CampaignConditions:  conditions,
		DateUpdate:          dateUpdate,
	}
}
