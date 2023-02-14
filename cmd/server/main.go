package main

import (
	"Comagic/internal/app_server"
	"Comagic/internal/config"
	"context"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/rs/zerolog"
	"os"
	"runtime/debug"
	"time"
)

const appName = "Comagic (Service)"

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
	cfg, err := config.NewServerConfig(".env")
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка в файле настроек")
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

	err = notify.Send(context.Background(), appName, "Скрипт запущен")
	if err != nil {
		logger.Warn().Err(err).Msg("ошибка сервера уведомлений")
	}
	a := app_server.NewApp(cfg, &logger)

	ctx := context.Background()
	err = a.Run(ctx)
	if err != nil {
		logger.Err(err).Msg("Ошибка выполнения запроса")
	}
}
