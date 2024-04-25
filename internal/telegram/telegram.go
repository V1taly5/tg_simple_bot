package telegram

import (
	"context"
	"log/slog"
	handlers "simple_tg_bot/internal/telegram/handlers"
	"simple_tg_bot/internal/telegram/mux"
	"simple_tg_bot/internal/yandex_disk/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgClient struct {
	api *tgbotapi.BotAPI
}

func NewTGClient(api *tgbotapi.BotAPI) *TgClient {
	return &TgClient{
		api: api,
	}
}

func (t *TgClient) RunTgClient(ctx context.Context, log *slog.Logger, uc *usecase.DiskUseCase) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.api.GetUpdatesChan(u)

	mux := mux.NewMux()

	mux.AddHandler(handlers.GetDisk())
	mux.AddHandler(handlers.DocHandl(log))
	mux.AddHandler(handlers.UploadFileCBHandl(ctx, log, uc))

	log.Info("[TG-BOT] launched")
	for {
		select {
		case update := <-updates:
			mux.Dispatch(t.api, update)
		case <-ctx.Done():
			log.Info("[TG-BOT] Stoped")
			return ctx.Err()
		}
	}
}
