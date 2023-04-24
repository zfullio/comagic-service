package service

import (
	"Comagic/internal/domain/entity"
	"Comagic/pkg/csv"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type CallRepositoryTracking interface {
	GetByDate(dateFrom time.Time, dateTill time.Time, fields []string) (calls []entity.Call, err error)
}

type CallRepositoryBQ interface {
	SendFromCS(ctx context.Context, bucket string, object string) (err error)
	DeleteByDateColumn(ctx context.Context, dateColumn string, dateStart time.Time, dateFinish time.Time) (err error)
	CreateTable(ctx context.Context) (err error)
}

type CallRepositoryCS interface {
	SendFile(ctx context.Context, filename string) (err error)
}

type CallService struct {
	tracking CallRepositoryTracking
	bq       CallRepositoryBQ
	cs       CallRepositoryCS
	logger   *zerolog.Logger
}

func NewCallService(tracking CallRepositoryTracking, bq CallRepositoryBQ, cs CallRepositoryCS, logger *zerolog.Logger) *CallService {
	serviceLogger := logger.With().Str("service", "call").Logger()

	return &CallService{
		tracking: tracking,
		bq:       bq,
		cs:       cs,
		logger:   &serviceLogger,
	}
}

func (s CallService) GetByDate(dateFrom time.Time, dateTill time.Time, fields []string) (calls []entity.Call, err error) {
	s.logger.Trace().Msg("GetByDate")
	calls, err = s.tracking.GetByDate(dateFrom, dateTill, fields)
	if err != nil {
		return calls, err
	}

	return calls, err
}

func (s CallService) SendAll(ctx context.Context, dateFrom time.Time, dateTill time.Time, bucketName string, calls []entity.Call) (err error) {
	s.logger.Trace().Msg("SendAll")

	dataForSend := make([]entity.CallCSV, 0, len(calls))
	for _, call := range calls {
		item := NewCallCSV(call)
		dataForSend = append(dataForSend, *item)
	}

	filename, err := csv.GenerateFile(dataForSend, "comagic_calls")
	if err != nil {
		return fmt.Errorf("ошибка генерации csv файла: %w", err)
	}

	err = s.cs.SendFile(ctx, filename)
	if err != nil {
		return fmt.Errorf("ошибка заливки на storage: %w", err)
	}

	err = s.bq.CreateTable(ctx)
	if err != nil {
		return fmt.Errorf("ошибка создания bq таблицы: %w", err)
	}

	s.logger.Info().Msgf("Удаление за %s -- %s", dateFrom.Format(time.DateOnly), dateTill.Format(time.DateOnly))

	err = s.bq.DeleteByDateColumn(ctx, "date", dateFrom, dateTill)
	if err != nil {
		return fmt.Errorf("ошибка удаления из bq: %w", err)
	}

	err = s.bq.SendFromCS(ctx, bucketName, filename)
	if err != nil {
		return fmt.Errorf("ошибка добавления в bq из storage: %w", err)
	}

	return nil
}

func NewCallCSV(call entity.Call) *entity.CallCSV {
	return &entity.CallCSV{
		Date:                          strings.Split(call.StartTime, " ")[0],
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
		LastAnsweredEmployeeRating:    call.LastAnsweredEmployeeId,
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
		OpenstatAd:                    call.OpenStatAd,
		OpenstatCampaign:              call.OpenStatCampaign,
		OpenstatService:               call.OpenStatService,
		OpenstatSource:                call.OpenStatSource,
		EqUtmSource:                   call.EqUtmSource,
		EqUtmMedium:                   call.EqUtmMedium,
		EqUtmTerm:                     call.EqUtmTerm,
		EqUtmContent:                  call.EqUtmContent,
		EqUtmCampaign:                 call.EqUtmCampaign,
		EqUtmReferrer:                 call.EqUtmReferrer,
		EqUtmExpid:                    call.EqUtmExpid,
		Attributes:                    call.Attributes,
		Tags:                          call.Tags,
		DateUpdate:                    call.DateUpdate.Format("2006-01-02 15:04:05.000"),
	}
}

func (s CallService) PushCallsToBQ(dateFrom time.Time, dateTill time.Time, fields []string, bucketName string) (err error) {
	s.logger.Trace().Msg("PushCallsToBQ")

	calls, err := s.GetByDate(dateFrom, dateTill, fields)
	if err != nil {
		return err
	}

	if len(calls) == 0 {
		return fmt.Errorf("звонки | пустой список значений")
	}

	ctx := context.Background()
	err = s.SendAll(ctx, dateFrom, dateTill, bucketName, calls)
	if err != nil {
		return fmt.Errorf("ошибка отправки звонков: %w", err)
	}

	return err
}
