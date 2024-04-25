package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"simple_tg_bot/internal/telegram/mux"
	"simple_tg_bot/internal/yandex_disk/usecase"

	yrepository "simple_tg_bot/internal/yandex_disk/y_repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	InfoMSG     = "сохраняется..."
	InfoDiskMSG = "запрос отправлен"
)

func DocHandl(log *slog.Logger) *mux.Handler {
	handl := func(u *mux.Update) {
		log.Info(fmt.Sprintf("got a new update(document): [from]: %s [subject]: %s", u.Message.From.UserName, u.Message.Document.FileName))

		saveInlineButton := tgbotapi.NewInlineKeyboardButtonData("Сохранить", "SaveFile")
		inlineRow := tgbotapi.NewInlineKeyboardRow(saveInlineButton)
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(inlineRow)

		msg := tgbotapi.NewMessage(u.Message.From.ID, "Получен файл. Что с ним делать?")
		msg.ReplyMarkup = inlineKeyboard

		msg.ReplyToMessageID = u.Message.MessageID

		_, err := u.Bot.Send(msg)
		if err != nil {
			log.Info(err.Error())
		}
	}
	handls := []mux.HandleFunc{handl}
	return mux.NewHandler(mux.And(mux.IsMessage(), mux.HasDocument()), handls)
}

func UploadFileCBHandl(ctx context.Context, log *slog.Logger, usecase *usecase.DiskUseCase) *mux.Handler {
	handl := func(u *mux.Update) {
		log.Info(fmt.Sprintf("got a new update(callback): [from]: %s [subject]: %s", u.CallbackQuery.Message.ReplyToMessage.From.UserName, u.CallbackQuery.Data))

		msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, InfoMSG)
		if _, err := u.Bot.Send(msg); err != nil {
			log.Warn(err.Error())
		}

		pathToSave := u.CallbackQuery.Message.ReplyToMessage.Caption
		fields := []string{"href", "method", "templated"}

		fileNameToSave := u.CallbackQuery.Message.ReplyToMessage.Document.FileName

		fileDownloadURL, err := u.Bot.GetFileDirectURL(u.CallbackQuery.Message.ReplyToMessage.Document.FileID)
		if err != nil {
			log.Error(err.Error())
		}

		// usecase
		_, err = usecase.UploudFileByURL(ctx, log, pathToSave, fileNameToSave, fileDownloadURL, true, fields)
		if err != nil {
			if errors.Is(err, yrepository.ErrResourceNotFound) {
				msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "Такого пути не существует")
				if _, err := u.Bot.Send(msg); err != nil {
					log.Warn(err.Error())
				}
			}
			log.Error(err.Error())
			return
		}

		msg1 := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.ReplyToMessage.Document.FileID)
		if _, err := u.Bot.Send(msg1); err != nil {
			log.Warn(err.Error())
		}

		deleteCallBackMSG := tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
		if resp, err := u.Bot.Request(deleteCallBackMSG); err != nil || !resp.Ok {
			log.Warn(fmt.Sprintf("failed to delete mesage id %d (%s): %v", deleteCallBackMSG.MessageID, string(resp.Result), err))
		}
	}
	// handls := []mux.HandleFunc{handl}
	return mux.NewCallBackQuertHandler("SaveFile", mux.IsCallBackQuery(), handl)
}

func GetDisk() *mux.Handler {
	handl := func(u *mux.Update) {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "sdfsdfsdf")

		u.Bot.Send(msg)
	}

	ll := mux.NewCommandHandler("help", mux.IsMessage(), handl)
	return ll
}
