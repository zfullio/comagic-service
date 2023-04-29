package policy

import (
	"Comagic/internal/domain/service"
	"time"
)

type CallPolicy struct {
	Service service.CallService
}

func NewCallPolicy(service service.CallService) *CallPolicy {
	return &CallPolicy{Service: service}
}

func (cp CallPolicy) PushCallsToBQ(dateFrom time.Time, dateTill time.Time, fields []string, bucketName string) (err error) {
	dateFromOnlyDate := time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, time.UTC)
	dateTillOnlyDate := time.Date(dateTill.Year(), dateTill.Month(), dateTill.Day(), 0, 0, 0, 0, time.UTC)

	err = cp.Service.PushCallsToBQ(dateFromOnlyDate, dateTillOnlyDate, fields, bucketName)
	if err != nil {
		return err
	}

	return err
}
