package comagic

import (
	"Comagic/internal/domain/entity"
	cm "Comagic/pkg/comagic"
	"fmt"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type callRepository struct {
	client *cm.Client
	logger *zerolog.Logger
}

func NewCallRepository(tracking *cm.Client, logger *zerolog.Logger) *callRepository {
	cmLogger := logger.With().Str("repo", "call").Str("type", "comagic").Logger()

	return &callRepository{
		client: tracking,
		logger: &cmLogger,
	}
}

func (cr callRepository) GetByDate(dateFrom time.Time, dateTill time.Time, fields []string) (calls []entity.Call, err error) {
	cr.logger.Trace().Msgf("GetByDate: %v -- %v", dateFrom.Format(time.DateOnly), dateTill.Format(time.DateOnly))

	callsFromRepo, err := cr.client.GetCallsReport(dateFrom, dateTill, fields)
	if err != nil {
		return calls, fmt.Errorf("ошибка получения звонков: %w", err)
	}

	t := time.Now()
	for i := 0; i < len(callsFromRepo); i++ {
		item := newCall(callsFromRepo[i], t)
		calls = append(calls, *item)
	}

	return calls, nil
}

func newCall(call cm.CallInfo, dateUpdate time.Time) *entity.Call {
	tagNames := make([]string, 0, len(call.Tags))
	for _, tag := range call.Tags {
		tagNames = append(tagNames, tag.TagName)
	}

	return &entity.Call{
		ID:                            call.ID,
		StartTime:                     call.StartTime,
		FinishTime:                    call.FinishTime,
		VirtualPhoneNumber:            call.VirtualPhoneNumber,
		IsTransfer:                    call.IsTransfer,
		FinishReason:                  call.FinishReason,
		Direction:                     call.Direction,
		Source:                        call.Source,
		CommunicationNumber:           call.CommunicationNumber,
		CommunicationPageURL:          call.CommunicationPageURL,
		CommunicationID:               call.CommunicationID,
		CommunicationType:             call.CommunicationType,
		IsLost:                        call.IsLost,
		CpnRegionID:                   call.CpnRegionID,
		CpnRegionName:                 call.CpnRegionName,
		WaitDuration:                  call.WaitDuration,
		TotalWaitDuration:             call.TotalWaitDuration,
		LostCallProcessingDuration:    call.LostCallProcessingDuration,
		TalkDuration:                  call.TalkDuration,
		CleanTalkDuration:             call.CleanTalkDuration,
		TotalDuration:                 call.TotalDuration,
		PostprocessDuration:           call.PostprocessDuration,
		UaClientID:                    call.UaClientID,
		YmClientID:                    call.YmClientID,
		SaleDate:                      call.SaleDate,
		SaleCost:                      call.SaleCost,
		SearchQuery:                   call.SearchQuery,
		SearchEngine:                  call.SearchEngine,
		ReferrerDomain:                call.ReferrerDomain,
		Referrer:                      call.Referrer,
		EntrancePage:                  call.EntrancePage,
		Gclid:                         call.Gclid,
		Yclid:                         call.Yclid,
		Ymclid:                        call.Ymclid,
		EfID:                          call.EfID,
		Channel:                       call.Channel,
		SiteID:                        call.SiteID,
		SiteDomainName:                call.SiteDomainName,
		CampaignID:                    call.CampaignID,
		CampaignName:                  call.CampaignName,
		AutoCallCampaignName:          call.AutoCallCampaignName,
		VisitOtherCampaign:            call.VisitOtherCampaign,
		VisitorID:                     call.VisitorID,
		PersonID:                      call.PersonID,
		VisitorType:                   call.VisitorType,
		VisitorSessionID:              call.VisitorSessionID,
		VisitsCount:                   call.VisitsCount,
		VisitorFirstCampaignID:        call.VisitorFirstCampaignID,
		VisitorFirstCampaignName:      call.VisitorFirstCampaignName,
		VisitorCity:                   call.VisitorCity,
		VisitorRegion:                 call.VisitorRegion,
		VisitorCountry:                call.VisitorCountry,
		VisitorDevice:                 call.VisitorDevice,
		LastAnsweredEmployeeID:        call.LastAnsweredEmployeeID,
		LastAnsweredEmployeeFullName:  call.LastAnsweredEmployeeFullName,
		LastAnsweredEmployeeRating:    call.LastAnsweredEmployeeRating,
		FirstAnsweredEmployeeID:       call.FirstAnsweredEmployeeID,
		FirstAnsweredEmployeeFullName: call.FirstAnsweredEmployeeFullName,
		ScenarioID:                    call.ScenarioID,
		ScenarioName:                  call.ScenarioName,
		CallAPIExternalID:             call.CallAPIExternalID,
		CallAPIRequestID:              call.CallAPIRequestID,
		ContactPhoneNumber:            call.ContactPhoneNumber,
		ContactFullName:               call.ContactFullName,
		ContactID:                     call.ContactID,
		UtmSource:                     call.UtmSource,
		UtmMedium:                     call.UtmMedium,
		UtmTerm:                       call.UtmTerm,
		UtmContent:                    call.UtmContent,
		UtmCampaign:                   call.UtmCampaign,
		OpenStatAd:                    call.OpenstatAd,
		OpenStatCampaign:              call.OpenstatCampaign,
		OpenStatService:               call.OpenstatService,
		OpenStatSource:                call.OpenstatSource,
		EqUtmSource:                   call.EqUtmSource,
		EqUtmMedium:                   call.EqUtmMedium,
		EqUtmTerm:                     call.EqUtmTerm,
		EqUtmContent:                  call.EqUtmContent,
		EqUtmCampaign:                 call.EqUtmCampaign,
		EqUtmReferrer:                 call.EqUtmReferrer,
		EqUtmExpid:                    call.EqUtmExpid,
		Attributes:                    strings.Join(call.Attributes, ", "),
		Tags:                          strings.Join(tagNames, ", "),
		DateUpdate:                    dateUpdate,
	}
}
