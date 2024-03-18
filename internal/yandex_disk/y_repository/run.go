package yrepository

import (
	"log/slog"
	"net/http"
)

const (
	baseUrl = "https://cloud-api.yandex.net/v1"
)

func NewYaDisk(log *slog.Logger, tokenStr string) *YandexDisk {
	token := new(Token)

	token.AccessToken = tokenStr

	httpClient := http.DefaultClient

	yDisk, err := NewYandexDisk(httpClient, token, baseUrl)
	if err != nil {
		log.Error(err.Error())
	}
	return yDisk
}
