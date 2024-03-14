package cmds

import (
	"context"
	"simple_tg_bot/internal/telegram/models"
	yandex "simple_tg_bot/internal/yandex_disk/y_repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	InfoDiskMSG = "запрос отправлен"
)

func CmdGetDisk(y *yandex.YandexDisk) models.CmdFunc {
	return func(ctx context.Context, api *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, InfoDiskMSG)
		if _, err := api.Send(msg); err != nil {
			return err
		}

		fieldMas := make([]string, 2)
		// fieldMas[0] = "paid_max_file_size"
		// fieldMas[1] = "user"

		disk, err := y.GetDisk(ctx, fieldMas)
		if err != nil {
			return err
		}
		msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, disk.User.Uid)
		if _, err := api.Send(msg1); err != nil {
			return err
		}
		return nil
	}
}
