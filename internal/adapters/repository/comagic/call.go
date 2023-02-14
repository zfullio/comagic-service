package comagic

import (
	"Comagic/internal/domain/entity"
	cm "Comagic/pkg/comagic"
	"fmt"
	"strings"
	"time"
)

type callRepository struct {
	client cm.Client
}

func NewCallRepository(tracking cm.Client) *callRepository {
	return &callRepository{client: tracking}
}

func (cr callRepository) GetByDate(dateFrom time.Time, dateTill time.Time, fields []string) (calls []entity.Call, err error) {
	callsFromRepo, err := cr.client.GetCallsReport(dateFrom, dateTill, fields)
	if err != nil {
		return calls, fmt.Errorf("ошибка получения получения звонков: %w", err)
	}
	t := time.Now()
	for i := 0; i < len(callsFromRepo); i++ {
		item := newCall(callsFromRepo[i], t)
		calls = append(calls, *item)
	}

	return calls, err
}

func newCall(call cm.CallInfo, dateUpdate time.Time) *entity.Call {
	tagNames := make([]string, 0, len(call.Tags))
	for _, tag := range call.Tags {
		tagNames = append(tagNames, tag.TagName)
	}
	return &entity.Call{
		Id:                            call.Id,
		StartTime:                     call.StartTime,
		FinishTime:                    call.FinishTime,
		VirtualPhoneNumber:            call.VirtualPhoneNumber,
		IsTransfer:                    call.IsTransfer,
		FinishReason:                  call.FinishReason,
		Direction:                     call.Direction,
		Source:                        call.Source,
		CommunicationNumber:           call.CommunicationNumber,
		CommunicationPageUrl:          call.CommunicationPageUrl,
		CommunicationId:               call.CommunicationId,
		CommunicationType:             call.CommunicationType,
		IsLost:                        call.IsLost,
		CpnRegionId:                   call.CpnRegionId,
		CpnRegionName:                 call.CpnRegionName,
		WaitDuration:                  call.WaitDuration,
		TotalWaitDuration:             call.TotalWaitDuration,
		LostCallProcessingDuration:    call.LostCallProcessingDuration,
		TalkDuration:                  call.TalkDuration,
		CleanTalkDuration:             call.CleanTalkDuration,
		TotalDuration:                 call.TotalDuration,
		PostprocessDuration:           call.PostprocessDuration,
		UaClientId:                    call.UaClientId,
		YmClientId:                    call.YmClientId,
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
		EfId:                          call.EfId,
		Channel:                       call.Channel,
		SiteId:                        call.SiteId,
		SiteDomainName:                call.SiteDomainName,
		CampaignId:                    call.CampaignId,
		CampaignName:                  call.CampaignName,
		AutoCallCampaignName:          call.AutoCallCampaignName,
		VisitOtherCampaign:            call.VisitOtherCampaign,
		VisitorId:                     call.VisitorId,
		PersonId:                      call.PersonId,
		VisitorType:                   call.VisitorType,
		VisitorSessionId:              call.VisitorSessionId,
		VisitsCount:                   call.VisitsCount,
		VisitorFirstCampaignId:        call.VisitorFirstCampaignId,
		VisitorFirstCampaignName:      call.VisitorFirstCampaignName,
		VisitorCity:                   call.VisitorCity,
		VisitorRegion:                 call.VisitorRegion,
		VisitorCountry:                call.VisitorCountry,
		VisitorDevice:                 call.VisitorDevice,
		LastAnsweredEmployeeId:        call.LastAnsweredEmployeeId,
		LastAnsweredEmployeeFullName:  call.LastAnsweredEmployeeFullName,
		LastAnsweredEmployeeRating:    call.LastAnsweredEmployeeRating,
		FirstAnsweredEmployeeId:       call.FirstAnsweredEmployeeId,
		FirstAnsweredEmployeeFullName: call.FirstAnsweredEmployeeFullName,
		ScenarioId:                    call.ScenarioId,
		ScenarioName:                  call.ScenarioName,
		CallApiExternalId:             call.CallApiExternalId,
		CallApiRequestId:              call.CallApiRequestId,
		ContactPhoneNumber:            call.ContactPhoneNumber,
		ContactFullName:               call.ContactFullName,
		ContactId:                     call.ContactId,
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
