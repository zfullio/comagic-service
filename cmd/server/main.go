package main

import (
	"Comagic/internal/app/server"
	"Comagic/internal/config"
	"context"
	"flag"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/rs/zerolog"
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

	var logger zerolog.Logger
	if *trace {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Caller().
			Int("pid", os.Getpid()).
			Str("go_version", buildInfo.GoVersion).
			Logger()
		logger.Info().Msg("Logging level = Trace")
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	if !*useEnv {
		logger.Info().Msgf("configuration file: %s", *fileConfig)
	} else {
		logger.Info().Msg("configuration from ENV")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *trace {
		logger.Info().Msg("Logging level = Trace")
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	cfg, err := config.NewServerConfig(*fileConfig, *useEnv)
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

	a := server.NewApp(cfg, &logger, appNotify)

	ctx := context.Background()
	err = a.Run(ctx)
	if err != nil {
		logger.Err(err).Msg("Ошибка выполнения запроса")
	}
}
