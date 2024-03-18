package app

import (
	"context"
	"log/slog"
	"os"
	"simple_tg_bot/internal/lib/logger/sl"
	"simple_tg_bot/internal/telegram"
	"simple_tg_bot/internal/telegram/menu"
	"simple_tg_bot/internal/telegram/menu/cmds"
	yandex "simple_tg_bot/internal/yandex_disk/y_repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	log    *slog.Logger
	menu   *menu.Menu
	client *telegram.TgClient
	yDisk  *yandex.YandexDisk
}

func New(bot *tgbotapi.BotAPI, yToken string, log *slog.Logger) *App {
	app := &App{}

	app.log = log

	app.yDisk = yandex.NewYaDisk(app.log, yToken)

	app.menu = menu.New()
	app.menu.RegisterCmd("help", cmds.CmdHelp())
	app.menu.RegisterCmd("getDisk", cmds.CmdGetDisk(app.yDisk))

	app.client = telegram.New(bot, app.menu)

	return app
}

func (a App) Run(ctx context.Context) error {
	const op = "app.run"
	log := a.log.With("op", op)
	err := a.client.RunTgClient(ctx, a.log)
	if err != nil {
		log.Error("error when starting the app", sl.Err(err))
		os.Exit(1)
	}
	return nil
}
