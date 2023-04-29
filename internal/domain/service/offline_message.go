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

type OfflineMessagesRepositoryTracking interface {
	GetByDate(dateFrom time.Time, dateTill time.Time, fields []string) (messages []entity.OfflineMessage, err error)
}

type OfflineMessageRepositoryBQ interface {
	SendFromCS(ctx context.Context, bucket string, object string) (err error)
	DeleteByDateColumn(ctx context.Context, dateStart time.Time, dateFinish time.Time) (err error)
	CreateTable(ctx context.Context) (err error)
	TableExists(ctx context.Context) (err error)
}

type OfflineMessageRepositoryCS interface {
	SendFile(ctx context.Context, filename string) (err error)
}

type OfflineMessageService struct {
	tracking OfflineMessagesRepositoryTracking
	bq       OfflineMessageRepositoryBQ
	cs       OfflineMessageRepositoryCS
	logger   *zerolog.Logger
}

func NewOfflineMessageService(tracking OfflineMessagesRepositoryTracking, bq OfflineMessageRepositoryBQ, cs OfflineMessageRepositoryCS, logger *zerolog.Logger) *OfflineMessageService {
	serviceLogger := logger.With().Str("service", "offline-message").Logger()

	return &OfflineMessageService{
		tracking: tracking,
		bq:       bq,
		cs:       cs,
		logger:   &serviceLogger,
	}
}

func (s OfflineMessageService) GetByDate(dateFrom time.Time, dateTill time.Time, fields []string) (messages []entity.OfflineMessage, err error) {
	s.logger.Trace().Msg("GetByDate")

	messages, err = s.tracking.GetByDate(dateFrom, dateTill, fields)
	if err != nil {
		return messages, err
	}

	return messages, err
}

func (s OfflineMessageService) SendAll(ctx context.Context, dateFrom time.Time, dateTill time.Time, bucketName string, messages []entity.OfflineMessage) (err error) {
	s.logger.Trace().Msg("SendAll")

	dataForSend := make([]entity.OfflineMessageCSV, 0, len(messages))

	for _, msg := range messages {
		item := NewOfflineMessageCSV(msg)
		dataForSend = append(dataForSend, *item)
	}

	filename, err := csv.GenerateFile(dataForSend, "comagic_offline_message")
	if err != nil {
		return fmt.Errorf("ошибка генерации csv файла: %w", err)
	}

	err = s.cs.SendFile(ctx, filename)
	if err != nil {
		return fmt.Errorf("ошибка заливки на storage: %w", err)
	}

	s.logger.Info().Msgf("Удаление за %s -- %s", dateFrom.Format(time.DateOnly), dateTill.Format(time.DateOnly))

	err = s.bq.DeleteByDateColumn(ctx, dateFrom, dateTill)
	if err != nil {
		return fmt.Errorf("ошибка удаления из bq: %w", err)
	}

	s.logger.Info().Msgf("Отправка файла в Cloud Storage: %s", filename)

	err = s.bq.SendFromCS(ctx, bucketName, filename)
	if err != nil {
		return fmt.Errorf("ошибка добавления в bq из storage: %w", err)
	}

	return err
}

func NewOfflineMessageCSV(message entity.OfflineMessage) *entity.OfflineMessageCSV {
	return &entity.OfflineMessageCSV{
		ID:                       message.ID,
		Date:                     strings.Split(message.DateTime, " ")[0],
		DateTime:                 message.DateTime,
		Text:                     message.Text,
		CommunicationNumber:      message.CommunicationNumber,
		CommunicationPageURL:     message.CommunicationPageURL,
		CommunicationType:        message.CommunicationType,
		CommunicationID:          message.CommunicationID,
		UaClientID:               message.UaClientID,
		YmClientID:               message.YmClientID,
		SaleDate:                 message.SaleDate,
		SaleCost:                 message.SaleCost,
		Status:                   message.Status,
		ProcessTime:              message.ProcessTime,
		FormType:                 message.FormType,
		FormName:                 message.FormName,
		SearchQuery:              message.SearchQuery,
		SearchEngine:             message.SearchEngine,
		ReferrerDomain:           message.ReferrerDomain,
		Referrer:                 message.Referrer,
		EntrancePage:             message.EntrancePage,
		Gclid:                    message.Gclid,
		Yclid:                    message.Yclid,
		Ymclid:                   message.Ymclid,
		EfID:                     message.EfID,
		Channel:                  message.Channel,
		EmployeeID:               message.EmployeeID,
		EmployeeFullName:         message.EmployeeFullName,
		EmployeeAnswerMessage:    message.EmployeeAnswerMessage,
		EmployeeComment:          message.EmployeeComment,
		Tags:                     message.Tags,
		SiteID:                   message.SiteID,
		SiteDomainName:           message.SiteDomainName,
		GroupID:                  message.GroupID,
		GroupName:                message.GroupName,
		CampaignID:               message.CampaignID,
		CampaignName:             message.CampaignName,
		VisitOtherCampaign:       message.VisitOtherCampaign,
		VisitorID:                message.VisitorID,
		VisitorName:              message.VisitorName,
		VisitorPhoneNumber:       message.VisitorPhoneNumber,
		VisitorEmail:             message.VisitorEmail,
		PersonID:                 message.PersonID,
		VisitorType:              message.VisitorType,
		VisitorSessionID:         message.VisitorSessionID,
		VisitsCount:              message.VisitsCount,
		VisitorFirstCampaignID:   message.VisitorFirstCampaignID,
		VisitorFirstCampaignName: message.VisitorFirstCampaignName,
		VisitorCity:              message.VisitorCity,
		VisitorRegion:            message.VisitorRegion,
		VisitorCountry:           message.VisitorCountry,
		VisitorDevice:            message.VisitorDevice,
		UtmSource:                message.UtmSource,
		UtmMedium:                message.UtmMedium,
		UtmTerm:                  message.UtmTerm,
		UtmContent:               message.UtmContent,
		UtmCampaign:              message.UtmCampaign,
		OpenStatAd:               message.OpenStatAd,
		OpenStatCampaign:         message.OpenStatCampaign,
		OpenStatService:          message.OpenStatService,
		OpenStatSource:           message.OpenStatSource,
		EqUtmSource:              message.EqUtmSource,
		EqUtmMedium:              message.EqUtmMedium,
		EqUtmTerm:                message.EqUtmTerm,
		EqUtmContent:             message.EqUtmContent,
		EqUtmCampaign:            message.EqUtmCampaign,
		EqUtmReferrer:            message.EqUtmReferrer,
		EqUtmExpid:               message.EqUtmExpid,
		Attributes:               message.Attributes,
		SourceID:                 message.SourceID,
		SourceName:               message.SourceName,
		SourceNew:                message.SourceNew,
		ChannelNew:               message.ChannelNew,
		ChannelCode:              message.ChannelCode,
		DateUpdate:               message.DateUpdate.Format("2006-01-02 15:04:05.000"),
	}
}

func (s OfflineMessageService) PushOfflineMessagesToBQ(dateFrom time.Time, dateTill time.Time, fields []string, bucketName string) (err error) {
	s.logger.Trace().Msg("PushOfflineMessagesToBQ")

	messages, err := s.GetByDate(dateFrom, dateTill, fields)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return fmt.Errorf("offline messages| пустой список значений")
	}

	ctx := context.Background()

	err = s.bq.TableExists(ctx)
	if err != nil {
		err = s.bq.CreateTable(ctx)
		if err != nil {
			return fmt.Errorf("ошибка создания bq таблицы: %w", err)
		}
	}

	err = s.SendAll(ctx, dateFrom, dateTill, bucketName, messages)
	if err != nil {
		return fmt.Errorf("ошибка отправки заявок: %w", err)
	}

	return nil
}
