package models

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CmdFunc func(ctx context.Context, api *tgbotapi.BotAPI, update *tgbotapi.Update) error
