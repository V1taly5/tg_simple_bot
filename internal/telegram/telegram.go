package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"simple_tg_bot/internal/lib/logger/sl"
	"simple_tg_bot/internal/telegram/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CmdProvider interface {
	RegisterCmd(cmd string, CmdFunc models.CmdFunc)
	Command(cmd string) (models.CmdFunc, bool)
}

type TgClient struct {
	api  *tgbotapi.BotAPI
	menu CmdProvider
}

func New(api *tgbotapi.BotAPI, menu CmdProvider) *TgClient {
	return &TgClient{
		api:  api,
		menu: menu,
	}
}

func (t *TgClient) HandleUpdate(ctx context.Context, update *tgbotapi.Update, log *slog.Logger) {
	const op = "telegram.HandlerUpdate"
	log = log.With(slog.String("op", op))

	switch {
	case isCommand(update):
		log.Info("got a new update: [from]: %s [subject]: %s", update.Message.From.UserName, update.Message.Text)

		cmd := update.Message.Command()
		cmdFunc, ok := t.menu.Command(cmd)
		if !ok {
			return
		}
		if err := cmdFunc(ctx, t.api, update); err != nil {
			log.Error("failed to handle update", sl.Err(err))
			if _, err := t.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s", "[Внутренняя ошибка]"))); err != nil {
				log.Error("failed to send message", sl.Err(err))
			}
		}
}

func (t *TgClient) RunTgClient(ctx context.Context, log *slog.Logger) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, 5*time.Second)
			t.HandleUpdate(updateCtx, &update, log)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func isCommand(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.IsCommand()
}
