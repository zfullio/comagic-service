package policy

import (
	"Comagic/internal/domain/service"
	"context"
	"fmt"
	"time"
)

type OfflineMessagePolicy struct {
	Service service.OfflineMessageService
}

func NewOfflineMessagePolicy(service service.OfflineMessageService) *OfflineMessagePolicy {
	return &OfflineMessagePolicy{Service: service}
}

func (cp OfflineMessagePolicy) PushOfflineMessageToBQ(ctx context.Context, dateFrom time.Time, dateTill time.Time, fields []string, bucketName string) error {
	dateFromOnlyDate := time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, time.UTC)
	dateTillOnlyDate := time.Date(dateTill.Year(), dateTill.Month(), dateTill.Day(), 0, 0, 0, 0, time.UTC)

	err := cp.Service.PushOfflineMessagesToBQ(ctx, dateFromOnlyDate, dateTillOnlyDate, fields, bucketName)
	if err != nil {
		return fmt.Errorf("service push offline messages to bq: %w", err)
	}

	return nil
}
