package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"simple_tg_bot/internal/lib/logger/sl"
	"simple_tg_bot/internal/telegram"
	"simple_tg_bot/internal/yandex_disk/usecase"
	yandex "simple_tg_bot/internal/yandex_disk/y_repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	log      *slog.Logger
	client   *telegram.TgClient
	yusecase *usecase.DiskUseCase
}

func NewApp(bot *tgbotapi.BotAPI, yToken string, log *slog.Logger) *App {
	app := &App{}
	app.log = log
	usecase := usecase.NewDiskUseCase(yandex.NewYaDisk(app.log, yToken))
	app.yusecase = usecase

	app.client = telegram.NewTGClient(bot)

	return app
}

func (a App) Run(ctx context.Context) error {
	const op = "app.run"
	log := a.log
	err := a.client.RunTgClient(ctx, log, a.yusecase)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		log.Error("error when starting the app", sl.Err(err))
		os.Exit(1)
	}
	return nil
}
