package cmds

import (
	"context"
	"simple_tg_bot/internal/telegram/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const HelpMSG = "This is help msg"

func CmdHelp() models.CmdFunc {
	return func(ctx context.Context, api *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, HelpMSG)
		if _, err := api.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
