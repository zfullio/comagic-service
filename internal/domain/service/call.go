package service

import (
	"Comagic/internal/domain/entity"
	"Comagic/pkg/csv"
	"context"
	"fmt"
	"github.com/dnlo/struct2csv"
	"github.com/rs/zerolog"
	"log"
	"os"
	"strings"
	"time"
)

type CallRepositoryTracking interface {
	GetByDate(dateFrom time.Time, dateTill time.Time, fields []string) (calls []entity.Call, err error)
}

type CallRepositoryBQ interface {
	SendFromCS(ctx context.Context, bucket string, object string) (err error)
	DeleteByDateColumn(ctx context.Context, dateColumn string, dateStart time.Time, dateFinish time.Time) (err error)
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
	s.logger.Info().Msg("GetByDate")
	calls, err = s.tracking.GetByDate(dateFrom, dateTill, fields)
	if err != nil {
		return calls, fmt.Errorf("ошибка получения получения звонков: %w", err)
	}
	return calls, err
}

func (s CallService) SendAll(ctx context.Context, dateFrom time.Time, dateTill time.Time, bucketName string, calls []entity.Call) (err error) {
	s.logger.Info().Msg("SendAll")
	dataForSend := make([]entity.CallCSV, 0, len(calls))
	for _, call := range calls {
		item := NewCallCSV(call)
		dataForSend = append(dataForSend, *item)
	}
	filename := csv.GenerateFilename("comagic_calls")
	err = csv.GenerateFile(dataForSend, filename)
	if err != nil {
		return fmt.Errorf("ошибка генерации csv файла: %w", err)
	}
	err = s.cs.SendFile(ctx, filename)
	if err != nil {
		return fmt.Errorf("ошибка заливки на storage: %w", err)
	}
	fmt.Printf("Удаление за %s -- %s", dateFrom, dateTill)
	err = s.bq.DeleteByDateColumn(ctx, "date", dateFrom, dateTill)
	if err != nil {
		return fmt.Errorf("ошибка удаления из bq: %w", err)
	}
	err = s.bq.SendFromCS(ctx, bucketName, filename)
	if err != nil {
		return fmt.Errorf("ошибка добавления в bq из storage: %w", err)
	}

	return err
}

func (s CallService) GenerateFile(calls []entity.CallCSV, name string) (err error) {
	s.logger.Info().Msg("GenerateFile")
	csvFile, err := os.Create(name)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()

	enc := struct2csv.New()
	_, err = enc.Marshal(calls)
	if err != nil {
		return err
	}
	w := struct2csv.NewWriter(csvFile)
	w.SetComma('|')
	err = w.WriteStructs(calls)
	if err != nil {
		return err
	}
	return err
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
	s.logger.Info().Msg("PushCallsToBQ")
	calls, err := s.GetByDate(dateFrom, dateTill, fields)
	if err != nil {
		return fmt.Errorf("ошибка получения звонков: %w", err)
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
