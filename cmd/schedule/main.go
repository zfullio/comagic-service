package main

import (
	"Comagic/internal/app_schedule"
	"Comagic/internal/config"
	"context"
	"flag"
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
	var fileConfig = flag.String("f", "config.yml", "configuration file")
	var useEnv = flag.Bool("env", false, "use environment variables")
	var trace = flag.Bool("trace", false, "switch trace logging")
	flag.Parse()

	buildInfo, _ := debug.ReadBuildInfo()
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	if !*useEnv {
		logger.Info().Msgf("configuration file: %s", *fileConfig)
	} else {
		logger.Info().Msg("configuration from ENV")
	}

	if *trace {
		logger.Level(zerolog.TraceLevel)
	} else {
		logger.Level(zerolog.InfoLevel)
	}

	cfg, err := config.NewScheduleConfig(*fileConfig, *useEnv)
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

	// Планировщик
	s := gocron.NewScheduler(time.UTC)
	s.WaitForScheduleAll()
	location, err := time.LoadLocation("Local")
	s.ChangeLocation(location)

	_, err = s.Every(1).Day().At(cfg.Time.All).Do(scheduleRun, a)
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
	dateFrom := dateTill.AddDate(0, -3, 0)

	err := a.PushCallsToBQ(dateFrom, dateTill)
	if err != nil {
		log.Printf("возникла ошибка во время PushCallsToBQ: %s", err)
		_ = a.Notify.Send(ctx, appName, fmt.Sprintf("возникла ошибка во время PushCallsToBQ: %s", err))
	}

	_ = a.Notify.Send(ctx, appName, fmt.Sprint("PushCallsToBQ: Успешно", err))

	err = a.PushOfflineMessagesToBQ(dateFrom, dateTill)
	if err != nil {
		log.Printf("возникла ошибка во время PushOfflineMessagesToBQ: %s", err)
		_ = a.Notify.Send(ctx, appName, fmt.Sprintf("возникла ошибка во время PushOfflineMessagesToBQ: %s", err))
	}

	_ = a.Notify.Send(ctx, appName, fmt.Sprint("PushOfflineMessagesToBQ: Успешно", err))

}
