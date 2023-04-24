package bq

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type offlineMessageRepository struct {
	db        bigquery.Client
	datasetID string
	tableID   string
	logger    *zerolog.Logger
}

func NewOfflineMessageRepository(client bigquery.Client, datasetID string, tableID string, logger *zerolog.Logger) *offlineMessageRepository {
	repoLogger := logger.With().Str("repo", "offline-message").Str("type", "bigquery").Logger()

	return &offlineMessageRepository{
		db:        client,
		datasetID: datasetID,
		tableID:   tableID,
		logger:    &repoLogger,
	}
}

func (or offlineMessageRepository) DeleteByDateColumn(ctx context.Context, dateColumn string, dateStart time.Time, dateFinish time.Time) (err error) {
	or.logger.Trace().Msgf("deleteByDateColumn: %v.%v", or.datasetID, or.tableID)

	err = DeleteByDateColumn(ctx, or.db, or.datasetID, or.tableID, dateColumn, dateStart, dateFinish)
	if err != nil {
		return fmt.Errorf("DeleteByDateColumn: %w", err)
	}

	return err
}

func (or offlineMessageRepository) CreateTable(ctx context.Context) (err error) {
	or.logger.Trace().Msgf("createTable: %v.%v", or.datasetID, or.tableID)

	myDataset := or.db.Dataset(or.datasetID)

	err = CreateTable(ctx, OfflineMessageDTO{}, myDataset, or.tableID)
	if err != nil {
		return fmt.Errorf("createTable: %w", err)
	}

	return nil
}

func (or offlineMessageRepository) SendFromCS(ctx context.Context, bucket string, object string) (err error) {
	or.logger.Trace().Msgf("sendFromCS: %v.%v", or.datasetID, or.tableID)

	myDataset := or.db.Dataset(or.datasetID)
	table := myDataset.Table(or.tableID)

	err = SendFromCS(ctx, OfflineMessageDTO{}, table, bucket, object)
	if err != nil {
		return fmt.Errorf("SendFromCS: %w", err)
	}

	return err
}

type OfflineMessageDTO struct {
	Id                       bigquery.NullInt64     `bigquery:"id"`
	Date                     bigquery.NullDate      `bigquery:"date"`
	DateTime                 bigquery.NullDateTime  `bigquery:"date_time"`
	Text                     bigquery.NullString    `bigquery:"text"`
	CommunicationNumber      bigquery.NullInt64     `bigquery:"communication_number"`
	CommunicationPageUrl     bigquery.NullString    `bigquery:"communication_page_url"`
	CommunicationType        bigquery.NullString    `bigquery:"communication_type"`
	CommunicationId          bigquery.NullInt64     `bigquery:"communication_id"`
	UaClientId               bigquery.NullString    `bigquery:"ua_client_id"`
	YmClientId               bigquery.NullString    `bigquery:"ym_client_id"`
	SaleDate                 bigquery.NullString    `bigquery:"sale_date"`
	SaleCost                 bigquery.NullFloat64   `bigquery:"sale_cost"`
	Status                   bigquery.NullString    `bigquery:"status"`
	ProcessTime              bigquery.NullString    `bigquery:"process_time"`
	FormType                 bigquery.NullString    `bigquery:"form_type"`
	FormName                 bigquery.NullString    `bigquery:"form_name"`
	SearchQuery              bigquery.NullString    `bigquery:"search_query"`
	SearchEngine             bigquery.NullString    `bigquery:"search_engine"`
	ReferrerDomain           bigquery.NullString    `bigquery:"referrer_domain"`
	Referrer                 bigquery.NullString    `bigquery:"referrer"`
	EntrancePage             bigquery.NullString    `bigquery:"entrance_page"`
	Gclid                    bigquery.NullString    `bigquery:"gclid"`
	Yclid                    bigquery.NullString    `bigquery:"yclid"`
	Ymclid                   bigquery.NullString    `bigquery:"ymclid"`
	EfId                     bigquery.NullString    `bigquery:"ef_id"`
	Channel                  bigquery.NullString    `bigquery:"channel"`
	EmployeeId               bigquery.NullInt64     `bigquery:"employee_id"`
	EmployeeFullName         bigquery.NullString    `bigquery:"employee_full_name"`
	EmployeeAnswerMessage    bigquery.NullString    `bigquery:"employee_answer_message"`
	EmployeeComment          bigquery.NullString    `bigquery:"employee_comment"`
	Tags                     bigquery.NullString    `bigquery:"tags"`
	SiteId                   bigquery.NullInt64     `bigquery:"site_id"`
	SiteDomainName           bigquery.NullString    `bigquery:"site_domain_name"`
	GroupId                  bigquery.NullInt64     `bigquery:"group_id"`
	GroupName                bigquery.NullString    `bigquery:"group_name"`
	CampaignId               bigquery.NullInt64     `bigquery:"campaign_id"`
	CampaignName             bigquery.NullString    `bigquery:"campaign_name"`
	VisitOtherCampaign       bigquery.NullBool      `bigquery:"visit_other_campaign"`
	VisitorId                bigquery.NullInt64     `bigquery:"visitor_id"`
	VisitorName              bigquery.NullString    `bigquery:"visitor_name"`
	VisitorPhoneNumber       bigquery.NullString    `bigquery:"visitor_phone_number"`
	VisitorEmail             bigquery.NullString    `bigquery:"visitor_email"`
	PersonId                 bigquery.NullInt64     `bigquery:"person_id"`
	VisitorType              bigquery.NullString    `bigquery:"visitor_type"`
	VisitorSessionId         bigquery.NullInt64     `bigquery:"visitor_session_id"`
	VisitsCount              bigquery.NullInt64     `bigquery:"visits_count"`
	VisitorFirstCampaignId   bigquery.NullInt64     `bigquery:"visitor_first_campaign_id"`
	VisitorFirstCampaignName bigquery.NullString    `bigquery:"visitor_first_campaign_name"`
	VisitorCity              bigquery.NullString    `bigquery:"visitor_city"`
	VisitorRegion            bigquery.NullString    `bigquery:"visitor_region"`
	VisitorCountry           bigquery.NullString    `bigquery:"visitor_country"`
	VisitorDevice            bigquery.NullString    `bigquery:"visitor_device"`
	UtmSource                bigquery.NullString    `bigquery:"utm_source"`
	UtmMedium                bigquery.NullString    `bigquery:"utm_medium"`
	UtmTerm                  bigquery.NullString    `bigquery:"utm_term"`
	UtmContent               bigquery.NullString    `bigquery:"utm_content"`
	UtmCampaign              bigquery.NullString    `bigquery:"utm_campaign"`
	OpenStatAd               bigquery.NullString    `bigquery:"openstat_ad"`
	OpenStatCampaign         bigquery.NullString    `bigquery:"openstat_campaign"`
	OpenStatService          bigquery.NullString    `bigquery:"openstat_service"`
	OpenStatSource           bigquery.NullString    `bigquery:"openstat_source"`
	EqUtmSource              bigquery.NullString    `bigquery:"eq_utm_source"`
	EqUtmMedium              bigquery.NullString    `bigquery:"eq_utm_medium"`
	EqUtmTerm                bigquery.NullString    `bigquery:"eq_utm_term"`
	EqUtmContent             bigquery.NullString    `bigquery:"eq_utm_content"`
	EqUtmCampaign            bigquery.NullString    `bigquery:"eq_utm_campaign"`
	EqUtmReferrer            bigquery.NullString    `bigquery:"eq_utm_referrer"`
	EqUtmExpid               bigquery.NullString    `bigquery:"eq_utm_expid"`
	Attributes               bigquery.NullString    `bigquery:"attributes"`
	SourceId                 bigquery.NullInt64     `bigquery:"source_id"`
	SourceName               bigquery.NullString    `bigquery:"source_name"`
	SourceNew                bigquery.NullString    `bigquery:"source_new"`
	ChannelNew               bigquery.NullString    `bigquery:"channel_new"`
	ChannelCode              bigquery.NullString    `bigquery:"channel_code"`
	DateUpdate               bigquery.NullTimestamp `bigquery:"date_update"`
}
