package main

import (
	"Comagic/internal/app_cli"
	"Comagic/internal/config"
	"context"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/rs/zerolog"
	_ "net/http/pprof"
	"os"
	"runtime/debug"
	"time"
)

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
	cfg, err := config.NewCliConfig("cli_config.yml")
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка в файле настроек")
	}

	comagicToken := "hwsxwp7a165ysxppexhegya71k03g3cqeyt9onn1"
	ctx := context.Background()

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

	a := app_cli.NewApp(ctx, cfg, comagicToken, &logger, appNotify)
	dateFrom := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	dateTill := time.Date(2023, 2, 10, 0, 0, 0, 0, time.UTC)
	err = a.PushCallsToBQ(dateFrom, dateTill)
	if err != nil {
		logger.Err(err).Msg("Ошибка выполнения запроса")
	}
}
