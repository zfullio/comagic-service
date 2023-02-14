package v1

import (
	cmBQ "Comagic/internal/adapters/repository/bq"
	cmRepo "Comagic/internal/adapters/repository/comagic"
	cmCS "Comagic/internal/adapters/repository/cs"
	"Comagic/internal/domain/policy"
	"Comagic/internal/domain/service"
	"Comagic/pb"
	"Comagic/pkg/comagic"
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/type/date"
	"time"
)

func (s Server) PushCallsToBQ(ctx context.Context, req *pb.PushCallsToBQRequest) (*pb.PushCallsToBQResponse, error) {
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
	clComagic := comagic.NewClient(comagic.DataAPI, s.cfg.Version, req.ComagicToken)
	cmCallRepo := cmRepo.NewCallRepository(*clComagic)

	bqClient, err := bigquery.NewClient(context.Background(), req.BqConfig.ProjectID, option.WithCredentialsFile(req.BqConfig.ServiceKey))
	if err != nil {
		return &pb.PushCallsToBQResponse{
			IsOK: false,
		}, fmt.Errorf("ошибка формирования клиента Big Query: %s", err)
	}
	bqCallRepo := cmBQ.NewCallRepository(*bqClient, req.BqConfig.DatasetID, req.BqConfig.TableID)

	csClient, err := storage.NewClient(ctx, option.WithCredentialsFile(req.CsConfig.ServiceKey))
	if err != nil {
		return &pb.PushCallsToBQResponse{
			IsOK: false,
		}, fmt.Errorf("ошибка формирования клиента Cloud Storage: %s", err)
	}
	csCallRepo := cmCS.NewCallRepository(*csClient, req.CsConfig.BucketName)

	srv := service.NewCallService(cmCallRepo, bqCallRepo, csCallRepo, s.logger)
	cmPolicy := policy.NewCallPolicy(*srv)

	dateFrom := pbDateNormalise(req.DateFrom)
	dateTill := pbDateNormalise(req.DateTill)
	err = cmPolicy.PushCallsToBQ(dateFrom, dateTill, fields, req.CsConfig.BucketName)
	if err != nil {
		return &pb.PushCallsToBQResponse{
			IsOK: false,
		}, fmt.Errorf("ошибка выполнения: %s", err)
	}

	return &pb.PushCallsToBQResponse{
		IsOK: true,
	}, nil
}

func pbDateNormalise(date *date.Date) time.Time {
	return time.Date(int(date.Year), time.Month(date.Month), int(date.Day), 0, 0, 0, 0, time.UTC)
}
