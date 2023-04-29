package comagic

import (
	"Comagic/internal/domain/entity"
	cm "Comagic/pkg/comagic"
	"fmt"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type offlineMessageRepository struct {
	client *cm.Client
	logger *zerolog.Logger
}

func NewOfflineMessageRepository(tracking *cm.Client, logger *zerolog.Logger) *offlineMessageRepository {
	cmLogger := logger.With().Str("repo", "offline-message").Str("type", "comagic").Logger()

	return &offlineMessageRepository{
		client: tracking,
		logger: &cmLogger,
	}
}

func (or offlineMessageRepository) GetByDate(dateFrom time.Time, dateTill time.Time, fields []string) (messages []entity.OfflineMessage, err error) {
	or.logger.Trace().Msgf("GetByDate: %v, %v", dateFrom, dateTill)

	messagesFromRepo, err := or.client.GetOfflineMessagesReport(dateFrom, dateTill, fields)
	if err != nil {
		return messages, fmt.Errorf("ошибка получения заявок: %w", err)
	}

	t := time.Now()
	for i := 0; i < len(messagesFromRepo); i++ {
		item := newOfflineMessage(messagesFromRepo[i], t)
		messages = append(messages, *item)
	}

	return messages, err
}

func newOfflineMessage(message cm.OfflineMessageInfo, dateUpdate time.Time) *entity.OfflineMessage {
	tagNames := make([]string, 0, len(message.Tags))
	for _, tag := range message.Tags {
		tagNames = append(tagNames, tag.TagName)
	}

	return &entity.OfflineMessage{
		ID:                       message.Id,
		DateTime:                 message.DateTime,
		Text:                     message.Text,
		CommunicationNumber:      message.CommunicationNumber,
		CommunicationPageURL:     message.CommunicationPageUrl,
		CommunicationType:        message.CommunicationType,
		CommunicationID:          message.CommunicationId,
		UaClientID:               message.UaClientId,
		YmClientID:               message.YmClientId,
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
		EfID:                     message.EfId,
		Channel:                  message.Channel,
		EmployeeID:               message.EmployeeId,
		EmployeeFullName:         message.EmployeeFullName,
		EmployeeAnswerMessage:    message.EmployeeAnswerMessage,
		EmployeeComment:          message.EmployeeComment,
		Tags:                     strings.Join(tagNames, ", "),
		SiteID:                   message.SiteId,
		SiteDomainName:           message.SiteDomainName,
		GroupID:                  message.GroupId,
		GroupName:                message.GroupName,
		CampaignID:               message.CampaignId,
		CampaignName:             message.CampaignName,
		VisitOtherCampaign:       message.VisitOtherCampaign,
		VisitorID:                message.VisitorId,
		VisitorName:              message.VisitorName,
		VisitorPhoneNumber:       message.VisitorPhoneNumber,
		VisitorEmail:             message.VisitorEmail,
		PersonID:                 message.PersonId,
		VisitorType:              message.VisitorType,
		VisitorSessionID:         message.VisitorSessionId,
		VisitsCount:              message.VisitsCount,
		VisitorFirstCampaignID:   message.VisitorFirstCampaignId,
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
		OpenStatAd:               message.OpenstatAd,
		OpenStatCampaign:         message.OpenstatCampaign,
		OpenStatService:          message.OpenstatService,
		OpenStatSource:           message.OpenstatSource,
		EqUtmSource:              message.EqUtmSource,
		EqUtmMedium:              message.EqUtmMedium,
		EqUtmTerm:                message.EqUtmTerm,
		EqUtmContent:             message.EqUtmContent,
		EqUtmCampaign:            message.EqUtmCampaign,
		EqUtmReferrer:            message.EqUtmReferrer,
		EqUtmExpid:               message.EqUtmExpid,
		Attributes:               strings.Join(message.Attributes, ", "),
		SourceID:                 message.SourceId,
		SourceName:               message.SourceName,
		SourceNew:                message.SourceNew,
		ChannelNew:               message.ChannelNew,
		ChannelCode:              message.ChannelCode,
		DateUpdate:               dateUpdate,
	}
}
