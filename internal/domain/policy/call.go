package policy

import (
	"Comagic/internal/domain/service"
	"context"
	"time"
)

type CallPolicy struct {
	Service service.CallService
}

func NewCallPolicy(service service.CallService) *CallPolicy {
	return &CallPolicy{Service: service}
}

func (cp CallPolicy) PushCallsToBQ(ctx context.Context, dateFrom time.Time, dateTill time.Time, fields []string, bucketName string) error {
	dateFromOnlyDate := time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, time.UTC)
	dateTillOnlyDate := time.Date(dateTill.Year(), dateTill.Month(), dateTill.Day(), 0, 0, 0, 0, time.UTC)

	err := cp.Service.PushCallsToBQ(ctx, dateFromOnlyDate, dateTillOnlyDate, fields, bucketName)
	if err != nil {
		return err
	}

	return nil
}
