package v1

import (
	"Comagic/internal/domain/policy"
	"Comagic/internal/domain/service"
	cmBQ "Comagic/internal/repository/bq"
	cmRepo "Comagic/internal/repository/comagic"
	cmCS "Comagic/internal/repository/cs"
	"Comagic/pb"
	"Comagic/pkg/comagic"
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
)

func (s Server) PushOfflineMessagesToBQ(ctx context.Context, req *pb.PushOfflineMessagesToBQRequest) (*pb.PushOfflineMessagesToBQResponse, error) {
	methodLogger := s.logger.With().Str("method", "PushOfflineMessagesToBQ").Str("client", req.BqConfig.ProjectId).Logger()

	methodLogger.Info().Msg(msgMethodPrepared)

	defer methodLogger.Info().Msg(msgMethodFinished)

	bqServiceKey := s.cfg.KeysDir + "/" + req.BqConfig.ServiceKey
	csServiceKey := s.cfg.KeysDir + "/" + req.CsConfig.ServiceKey

	dateFrom, err := pbDateNormalize(req.Period.DateFrom)
	if err != nil {
		methodLogger.Error().Err(err).Msg(msgErrMethod)

		return &pb.PushOfflineMessagesToBQResponse{
			IsOk: false,
		}, fmt.Errorf("wrong value in field 'dateFrom' : %w", err)
	}

	dateTill, err := pbDateNormalize(req.Period.DateTill)
	if err != nil {
		methodLogger.Error().Err(err).Msg(msgErrMethod)

		return &pb.PushOfflineMessagesToBQResponse{
			IsOk: false,
		}, fmt.Errorf("wrong value in field 'dateTill' : %w", err)
	}

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

	clComagic := comagic.NewClient(comagic.DataAPI, s.cfg.Version, req.ComagicToken)
	cmOfflineMessageRepo := cmRepo.NewOfflineMessageRepository(clComagic, &methodLogger)

	bqClient, err := bigquery.NewClient(ctx, req.BqConfig.ProjectId, option.WithCredentialsFile(bqServiceKey))
	if err != nil {
		methodLogger.Error().Err(err).Msg(msgErrMethod)

		return &pb.PushOfflineMessagesToBQResponse{
			IsOk: false,
		}, fmt.Errorf("ошибка формирования клиента Big Query: %s", err)
	}

	defer func(bqClient *bigquery.Client) {
		err := bqClient.Close()
		if err != nil {
			methodLogger.Err(err).Msg("ошибка закрытия клиента Big Query")
		}
	}(bqClient)

	bqOfflineMessageRepo := cmBQ.NewOfflineMessageRepository(bqClient, req.BqConfig.DatasetId, req.BqConfig.TableId, &methodLogger)

	csClient, err := storage.NewClient(ctx, option.WithCredentialsFile(csServiceKey))
	if err != nil {
		methodLogger.Error().Err(err).Msg(msgErrMethod)

		return &pb.PushOfflineMessagesToBQResponse{
			IsOk: false,
		}, fmt.Errorf("ошибка формирования клиента Cloud Storage: %s", err)
	}

	defer func(csClient *storage.Client) {
		err := csClient.Close()
		if err != nil {
			methodLogger.Err(err).Msg("ошибка закрытия клиента Cloud Storage")
		}
	}(csClient)

	csOfflineMessageRepo := cmCS.NewOfflineMessageRepository(csClient, req.CsConfig.BucketName, &methodLogger)

	srv := service.NewOfflineMessageService(cmOfflineMessageRepo, bqOfflineMessageRepo, csOfflineMessageRepo, &methodLogger)
	cmPolicy := policy.NewOfflineMessagePolicy(*srv)

	methodLogger.Info().Msg(msgMethodStarted)

	err = cmPolicy.PushOfflineMessageToBQ(ctx, dateFrom, dateTill.AddDate(0, 0, 1), fields, req.CsConfig.BucketName)
	if err != nil {
		methodLogger.Error().Err(err).Msg(msgErrMethod)

		return &pb.PushOfflineMessagesToBQResponse{
			IsOk: false,
		}, fmt.Errorf("ошибка выполнения: %s", err)
	}

	return &pb.PushOfflineMessagesToBQResponse{
		IsOk: true,
	}, nil
}
