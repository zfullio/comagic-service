package policy

import (
	"Comagic/internal/domain/service"
	"time"
)

type OfflineMessagePolicy struct {
	Service service.OfflineMessageService
}

func NewOfflineMessagePolicy(service service.OfflineMessageService) *OfflineMessagePolicy {
	return &OfflineMessagePolicy{Service: service}
}

func (cp OfflineMessagePolicy) PushOfflineMessageToBQ(dateFrom time.Time, dateTill time.Time, fields []string, bucketName string) (err error) {
	dateFromOnlyDate := time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, time.UTC)
	dateTillOnlyDate := time.Date(dateTill.Year(), dateTill.Month(), dateTill.Day(), 0, 0, 0, 0, time.UTC)

	err = cp.Service.PushOfflineMessagesToBQ(dateFromOnlyDate, dateTillOnlyDate, fields, bucketName)
	if err != nil {
		return err
	}

	return err
}
