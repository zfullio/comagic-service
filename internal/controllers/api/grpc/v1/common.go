package v1

import (
	"time"
)

const msgErrMethod = "ошибка выполнения"
const msgMethodPrepared = "Подготовка"
const msgMethodStarted = "Запущено"
const msgMethodFinished = "Завершено"

func pbDateNormalize(s string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
