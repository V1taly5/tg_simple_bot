package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"simple_tg_bot/internal/app"
	"simple_tg_bot/internal/config"
	"simple_tg_bot/internal/lib/logger/handlers/slogpretty"
	"simple_tg_bot/internal/lib/logger/sl"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	bot, err := tgbotapi.NewBotAPI(cfg.TgBotToken)
	if err != nil {
		log.Error("bot is not running", sl.Err(err))
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	App := app.New(bot, cfg.YandexDiskToken, log)
	App.Run(ctx)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

// togo:: App config: tgToken, yandexDiskTiken, env,
