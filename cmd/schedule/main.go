package main

import (
	"Comagic/internal/app_schedule"
	"Comagic/internal/config"
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/rs/zerolog"
	"log"
	"os"
	"runtime/debug"
	"time"
)

const appName = "Comagic (schedule)"

func main() {
	buildInfo, _ := debug.ReadBuildInfo()
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()
	cfg, err := config.NewScheduleConfig("schedule_config.yml")
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка в файле настроек")
	}
	err = cfg.Check()
	if err != nil {
		log.Fatalf("eror in config: %s", err)
	}
	telegramService, err := telegram.New(cfg.TG.Token)
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка в сервисе: Telegram")
	}
	telegramService.AddReceivers(cfg.Chat)

	appNotify := notify.New()
	appNotify.UseServices(telegramService)

	if !cfg.IsEnabled {
		notify.Disable(appNotify)
	}

	ctx := context.Background()
	a := app_schedule.NewApp(ctx, cfg, cfg.Comagic.Token, &logger, appNotify)

	s := gocron.NewScheduler(time.UTC)
	s.WaitForScheduleAll()
	location, err := time.LoadLocation("Local")
	s.ChangeLocation(location)

	_, err = s.Every(1).Day().At(cfg.Time.Calls).Do(scheduleRun, a)
	if err != nil {
		log.Fatalf("Ошибка планировщика %v", err)
	}
	s.StartAsync()
	_, t := s.NextRun()
	_ = appNotify.Send(ctx, appName, fmt.Sprintf("Next run: %s", t.Format(time.DateTime)))
	s.StartBlocking()
}

func scheduleRun(a *app_schedule.App) {
	ctx := context.Background()
	now := time.Now()
	dateTill := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dateFrom := dateTill.AddDate(0, -2, 0)
	log.Printf("Выполнение PushCallsToBQ за период: %s -- %s", dateFrom.Format(time.DateOnly), dateTill.Format(time.DateOnly))
	err := a.PushCallsToBQ(dateFrom, dateTill)
	if err != nil {
		log.Printf("возникла ошибка во время PushCallsToBQ: %s", err)
		_ = a.Notify.Send(ctx, appName, fmt.Sprintf("возникла ошибка во время PushCallsToBQ: %s", err))
	}
	_ = a.Notify.Send(ctx, appName, fmt.Sprint("PushCallsToBQ: Успешно", err))
}
