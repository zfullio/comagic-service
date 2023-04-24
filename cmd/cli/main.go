package main

import (
	"Comagic/internal/app_cli"
	"Comagic/internal/config"
	"context"
	"flag"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/rs/zerolog"
	_ "net/http/pprof"
	"os"
	"runtime/debug"
	"time"
)

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

	cfg, err := config.NewCliConfig(*fileConfig, *useEnv)
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка в файле настроек")
	}

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

	a := app_cli.NewApp(ctx, cfg, cfg.Comagic.Token, &logger, appNotify)
	dateFrom := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	dateTill := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	err = a.PushCallsToBQ(dateFrom, dateTill)
	if err != nil {
		logger.Err(err).Msg("Ошибка выполнения запроса")
	}
}
