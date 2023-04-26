package bq

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type callRepository struct {
	db     *bigquery.Client
	table  *bigquery.Table
	logger *zerolog.Logger
}

func NewCallRepository(client *bigquery.Client, datasetID string, tableID string, logger *zerolog.Logger) *callRepository {
	repoLogger := logger.With().Str("repo", "call").Str("type", "bigquery").Logger()

	dataset := client.Dataset(datasetID)
	table := dataset.Table(tableID)

	return &callRepository{
		db:     client,
		table:  table,
		logger: &repoLogger,
	}
}

func (cr callRepository) TableExists(ctx context.Context) (err error) {
	cr.logger.Trace().Msg("TableExists")

	err = TableExists(ctx, cr.table)
	if err != nil {
		return err
	}

	return nil
}

func (cr callRepository) CreateTable(ctx context.Context) (err error) {
	cr.logger.Trace().Msg("createTable")

	fieldClustering := []string{"campaign_id", "utm_source", "utm_medium"}
	fieldPartition := "date"
	err = CreateTable(ctx, CallDTO{}, cr.table, &fieldPartition, &fieldClustering)
	if err != nil {
		return fmt.Errorf("createTable: %w", err)
	}

	return nil
}

func (cr callRepository) DeleteByDateColumn(ctx context.Context, dateStart time.Time, dateFinish time.Time) (err error) {
	cr.logger.Trace().Msg("deleteByDateColumn")

	dateColumn := "date"

	err = DeleteByDateColumn(ctx, cr.db, cr.table, dateColumn, dateStart, dateFinish)
	if err != nil {
		return fmt.Errorf("DeleteByDateColumn: %w", err)
	}

	return err
}

func (cr callRepository) SendFromCS(ctx context.Context, bucket string, object string) (err error) {
	cr.logger.Trace().Msg("sendFromCS")

	err = SendFromCS(ctx, CallDTO{}, cr.table, bucket, object)
	if err != nil {
		return fmt.Errorf("SendFromCS: %w", err)
	}

	return err
}

type CallDTO struct {
	Date                          bigquery.NullDate      `bigquery:"date"`
	StartTime                     bigquery.NullDateTime  `bigquery:"start_time"`
	FinishTime                    bigquery.NullDateTime  `bigquery:"finish_time"`
	VirtualPhoneNumber            bigquery.NullString    `bigquery:"virtual_phone_number"`
	IsTransfer                    bigquery.NullBool      `bigquery:"is_transfer"`
	FinishReason                  bigquery.NullString    `bigquery:"finish_reason"`
	Direction                     bigquery.NullString    `bigquery:"direction"`
	Source                        bigquery.NullString    `bigquery:"source"`
	CommunicationNumber           bigquery.NullInt64     `bigquery:"communication_number"`
	CommunicationPageUrl          bigquery.NullString    `bigquery:"communication_page_url"`
	CommunicationId               bigquery.NullInt64     `bigquery:"communication_id"`
	CommunicationType             bigquery.NullString    `bigquery:"communication_type"`
	IsLost                        bigquery.NullBool      `bigquery:"is_lost"`
	CpnRegionId                   bigquery.NullInt64     `bigquery:"cpn_region_id"`
	CpnRegionName                 bigquery.NullString    `bigquery:"cpn_region_name"`
	WaitDuration                  bigquery.NullInt64     `bigquery:"wait_duration"`
	TotalWaitDuration             bigquery.NullInt64     `bigquery:"total_wait_duration"`
	LostCallProcessingDuration    bigquery.NullInt64     `bigquery:"lost_call_processing_duration"`
	TalkDuration                  bigquery.NullInt64     `bigquery:"talk_duration"`
	CleanTalkDuration             bigquery.NullInt64     `bigquery:"clean_talk_duration"`
	TotalDuration                 bigquery.NullInt64     `bigquery:"total_duration"`
	PostprocessDuration           bigquery.NullInt64     `bigquery:"postprocess_duration"`
	UaClientId                    bigquery.NullString    `bigquery:"ua_client_id"`
	YmClientId                    bigquery.NullString    `bigquery:"ym_client_id"`
	SaleDate                      bigquery.NullString    `bigquery:"sale_date"`
	SaleCost                      bigquery.NullFloat64   `bigquery:"sale_cost"`
	SearchQuery                   bigquery.NullString    `bigquery:"search_query"`
	SearchEngine                  bigquery.NullString    `bigquery:"search_engine"`
	ReferrerDomain                bigquery.NullString    `bigquery:"referrer_domain"`
	Referrer                      bigquery.NullString    `bigquery:"referrer"`
	EntrancePage                  bigquery.NullString    `bigquery:"entrance_page"`
	Gclid                         bigquery.NullString    `bigquery:"gclid"`
	Yclid                         bigquery.NullString    `bigquery:"yclid"`
	Ymclid                        bigquery.NullString    `bigquery:"ymclid"`
	EfId                          bigquery.NullString    `bigquery:"ef_id"`
	Channel                       bigquery.NullString    `bigquery:"channel"`
	SiteId                        bigquery.NullInt64     `bigquery:"site_id"`
	SiteDomainName                bigquery.NullString    `bigquery:"site_domain_name"`
	CampaignId                    bigquery.NullInt64     `bigquery:"campaign_id"`
	CampaignName                  bigquery.NullString    `bigquery:"campaign_name"`
	AutoCallCampaignName          bigquery.NullString    `bigquery:"auto_call_campaign_name"`
	VisitOtherCampaign            bigquery.NullBool      `bigquery:"visit_other_campaign"`
	VisitorId                     bigquery.NullInt64     `bigquery:"visitor_id"`
	PersonId                      bigquery.NullInt64     `bigquery:"person_id"`
	VisitorType                   bigquery.NullString    `bigquery:"visitor_type"`
	VisitorSessionId              bigquery.NullInt64     `bigquery:"visitor_session_id"`
	VisitsCount                   bigquery.NullInt64     `bigquery:"visits_count"`
	VisitorFirstCampaignId        bigquery.NullInt64     `bigquery:"visitor_first_campaign_id"`
	VisitorFirstCampaignName      bigquery.NullString    `bigquery:"visitor_first_campaign_name"`
	VisitorCity                   bigquery.NullString    `bigquery:"visitor_city"`
	VisitorRegion                 bigquery.NullString    `bigquery:"visitor_region"`
	VisitorCountry                bigquery.NullString    `bigquery:"visitor_country"`
	VisitorDevice                 bigquery.NullString    `bigquery:"visitor_device"`
	LastAnsweredEmployeeId        bigquery.NullString    `bigquery:"last_answered_employee_id"`
	LastAnsweredEmployeeFullName  bigquery.NullString    `bigquery:"last_answered_employee_full_name"`
	LastAnsweredEmployeeRating    bigquery.NullInt64     `bigquery:"last_answered_employee_rating"`
	FirstAnsweredEmployeeId       bigquery.NullInt64     `bigquery:"first_answered_employee_id"`
	FirstAnsweredEmployeeFullName bigquery.NullString    `bigquery:"first_answered_employee_full_name"`
	ScenarioId                    bigquery.NullString    `bigquery:"scenario_id"`
	ScenarioName                  bigquery.NullString    `bigquery:"scenario_name"`
	CallApiExternalId             bigquery.NullString    `bigquery:"call_api_external_id"`
	CallApiRequestId              bigquery.NullString    `bigquery:"call_api_request_id"`
	ContactPhoneNumber            bigquery.NullString    `bigquery:"contact_phone_number"`
	ContactFullName               bigquery.NullString    `bigquery:"contact_full_name"`
	ContactId                     bigquery.NullString    `bigquery:"contact_id"`
	UtmSource                     bigquery.NullString    `bigquery:"utm_source"`
	UtmMedium                     bigquery.NullString    `bigquery:"utm_medium"`
	UtmTerm                       bigquery.NullString    `bigquery:"utm_term"`
	UtmContent                    bigquery.NullString    `bigquery:"utm_content"`
	UtmCampaign                   bigquery.NullString    `bigquery:"utm_campaign"`
	OpenstatAd                    bigquery.NullString    `bigquery:"openstat_ad"`
	OpenstatCampaign              bigquery.NullString    `bigquery:"openstat_campaign"`
	OpenstatService               bigquery.NullString    `bigquery:"openstat_service"`
	OpenstatSource                bigquery.NullString    `bigquery:"openstat_source"`
	EqUtmSource                   bigquery.NullString    `bigquery:"eq_utm_source"`
	EqUtmMedium                   bigquery.NullString    `bigquery:"eq_utm_medium"`
	EqUtmTerm                     bigquery.NullString    `bigquery:"eq_utm_term"`
	EqUtmContent                  bigquery.NullString    `bigquery:"eq_utm_content"`
	EqUtmCampaign                 bigquery.NullString    `bigquery:"eq_utm_campaign"`
	EqUtmReferrer                 bigquery.NullString    `bigquery:"eq_utm_referrer"`
	EqUtmExpid                    bigquery.NullString    `bigquery:"eq_utm_expid"`
	Attributes                    bigquery.NullString    `bigquery:"attributes"`
	Tags                          bigquery.NullString    `bigquery:"tags"`
	DateUpdate                    bigquery.NullTimestamp `bigquery:"date_update"`
}
