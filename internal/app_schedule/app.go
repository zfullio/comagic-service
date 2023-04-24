package app_schedule

import (
	"Comagic/internal/adapters/repository/bq"
	"Comagic/internal/adapters/repository/comagic"
	"Comagic/internal/adapters/repository/cs"
	"Comagic/internal/config"
	"Comagic/internal/domain/service"
	cm "Comagic/pkg/comagic"
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/nikoksr/notify"
	"github.com/rs/zerolog"
	"google.golang.org/api/option"
	"log"
	"time"
)

type Entity int

type App struct {
	logger        *zerolog.Logger
	cfg           *config.ScheduleConfig
	callSrv       service.CallService
	offMessageSrv service.OfflineMessageService
	campaignSrv   service.CampaignService
	Notify        notify.Notifier
}

func NewApp(ctx context.Context, cfg *config.ScheduleConfig, token string, logger *zerolog.Logger, notify notify.Notifier) *App {
	clComagic := cm.NewClient(cm.DataAPI, cfg.Comagic.Version, token)
	cmCallRepo := comagic.NewCallRepository(*clComagic, logger)
	cmCampaignRepo := comagic.NewCampaignRepository(*clComagic, logger)
	cmOfflineMessageRepo := comagic.NewOfflineMessageRepository(*clComagic, logger)

	bqClient, err := bigquery.NewClient(context.Background(), cfg.BQ.ProjectID, option.WithCredentialsFile(cfg.BQ.ServiceKeyPath))
	if err != nil {
		log.Fatalf("ошибка получения клиента Big Query: %s", err)
	}
	bqCallRepo := bq.NewCallRepository(*bqClient, cfg.CallReport.DatasetID, cfg.CallReport.TableID, logger)
	bqCampaignRepo := bq.NewCampaignRepository(*bqClient, cfg.CampaignReport.DatasetID, cfg.CampaignReport.TableID, logger)
	bqOfflineMessageRepo := bq.NewOfflineMessageRepository(*bqClient, cfg.OfflineMessageReport.DatasetID, cfg.OfflineMessageReport.TableID, logger)

	csClient, err := storage.NewClient(ctx, option.WithCredentialsFile(cfg.CS.ServiceKeyPath))
	if err != nil {
		log.Fatalf("ошибка получения клиента Cloud Storage: %s", err)
	}
	csCallRepo := cs.NewCallRepository(*csClient, cfg.BucketName, logger)
	csOfflineMessageRepo := cs.NewOfflineMessageRepository(*csClient, cfg.BucketName, logger)

	callSrv := service.NewCallService(cmCallRepo, bqCallRepo, csCallRepo, logger)
	campaignSrv := service.NewCampaignService(cmCampaignRepo, bqCampaignRepo, logger)
	offlineMessagesSrv := service.NewOfflineMessageService(cmOfflineMessageRepo, bqOfflineMessageRepo, csOfflineMessageRepo, logger)

	return &App{
		logger:        logger,
		cfg:           cfg,
		callSrv:       *callSrv,
		campaignSrv:   *campaignSrv,
		offMessageSrv: *offlineMessagesSrv,
		Notify:        notify,
	}
}

func (a App) PushCallsToBQ(dateFrom time.Time, dateTill time.Time) (err error) {
	_ = a.Notify.Send(context.Background(), "Comagic", fmt.Sprintf("PushCallsToBQ: %s.%s.%s", a.cfg.BQ.ProjectID,
		a.cfg.CallReport.DatasetID, a.cfg.CallReport.TableID))
	a.logger.Info().Msgf("Get calls: %s -- %s", dateFrom.Format(time.DateOnly), dateTill.Format(time.DateOnly))

	fields := []string{"id", "start_time", "finish_time", "finish_reason", "direction", "cpn_region_id",
		"cpn_region_name", "scenario_operations", "scenario_id", "scenario_name", "source", "is_lost",
		"communication_number", "communication_page_url", "contact_phone_number", "communication_id", "communication_type",
		"wait_duration", "total_wait_duration", "lost_call_processing_duration", "talk_duration", "clean_talk_duration",
		"total_duration", "postprocess_duration", "call_records", "wav_call_records", "full_record_file_link",
		"voice_mail_records", "virtual_phone_number", "ua_client_id", "ym_client_id", "sale_date", "sale_cost", "is_transfer",
		"search_query", "search_engine", "referrer_domain", "referrer", "entrance_page", "gclid", "yclid", "ymclid", "ef_id",
		"channel", "tags", "employees", "last_answered_employee_id", "last_answered_employee_full_name", "last_answered_employee_rating",
		"first_answered_employee_id", "first_answered_employee_full_name", "last_talked_employee_id",
		"last_talked_employee_full_name", "first_talked_employee_id", "first_talked_employee_full_name", "scenario_name",
		"scenario_id", "site_domain_name", "site_id", "campaign_name", "campaign_id", "visit_other_campaign",
		"auto_call_campaign_name", "visitor_id", "person_id", "visitor_type", "visitor_session_id", "visits_count",
		"visitor_first_campaign_id", "visitor_first_campaign_name", "visitor_city", "visitor_region", "visitor_country",
		"visitor_device", "visitor_custom_properties", "segments", "call_api_request_id", "call_api_external_id", "contact_id",
		"contact_full_name", "utm_source", "utm_medium", "utm_term", "utm_content", "utm_campaign", "openstat_ad",
		"openstat_campaign", "openstat_service", "openstat_source", "attributes", "eq_utm_source", "eq_utm_medium",
		"eq_utm_term", "eq_utm_content", "eq_utm_campaign", "eq_utm_referrer", "eq_utm_expid",
	}

	err = a.callSrv.PushCallsToBQ(dateFrom, dateTill, fields, a.cfg.BucketName)
	if err != nil {
		return fmt.Errorf("не могу выполнить запрос: %s", err)
	}
	return nil
}

func (a App) PushOfflineMessagesToBQ(dateFrom time.Time, dateTill time.Time) (err error) {
	_ = a.Notify.Send(context.Background(), "Comagic", fmt.Sprintf("PushOfflineMessagesToBQ: %s.%s.%s", a.cfg.BQ.ProjectID,
		a.cfg.OfflineMessageReport.DatasetID, a.cfg.OfflineMessageReport.TableID))
	a.logger.Info().Msgf("Get offline-messages: %s -- %s", dateFrom.Format(time.DateOnly), dateTill.Format(time.DateOnly))

	fields := []string{
		"id",
		"date_time",
		"text",
		"communication_number",
		"communication_page_url",
		"communication_type",
		"communication_id",
		"ua_client_id",
		"ym_client_id",
		"sale_date",
		"sale_cost",
		"status",
		"process_time",
		"form_type",
		"form_name",
		"search_query",
		"search_engine",
		"referrer_domain",
		"referrer",
		"entrance_page",
		"gclid",
		"yclid",
		"ymclid",
		"ef_id",
		"channel",
		"employee_id",
		"employee_full_name",
		"employee_answer_message",
		"employee_comment",
		"tags",
		"site_id",
		"site_domain_name",
		"group_id",
		"group_name",
		"campaign_id",
		"campaign_name",
		"visit_other_campaign",
		"visitor_id",
		"visitor_name",
		"visitor_phone_number",
		"visitor_email",
		"person_id",
		"visitor_type",
		"visitor_session_id",
		"visits_count",
		"visitor_first_campaign_id",
		"visitor_first_campaign_name",
		"visitor_city",
		"visitor_region",
		"visitor_country",
		"visitor_device",
		"visitor_custom_properties",
		"segments",
		"utm_source",
		"utm_medium",
		"utm_term",
		"utm_content",
		"utm_campaign",
		"openstat_ad",
		"openstat_campaign",
		"openstat_service",
		"openstat_source",
		"eq_utm_source",
		"eq_utm_medium",
		"eq_utm_term",
		"eq_utm_content",
		"eq_utm_campaign",
		"eq_utm_referrer",
		"eq_utm_expid",
		"attributes",
		"source_id",
		"source_name",
		"source_new",
		"channel_new",
		"channel_code",
	}

	err = a.offMessageSrv.PushOfflineMessagesToBQ(dateFrom, dateTill, fields, a.cfg.BucketName)
	if err != nil {
		return fmt.Errorf("не могу выполнить запрос: %s", err)
	}
	return err
}
