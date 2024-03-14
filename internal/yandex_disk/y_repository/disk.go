package yrepository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func NewYandexDisk(client *http.Client, token *Token, baseUrl string) (*YandexDisk, error) {
	const op = "ydisk/NewYandexDisk"
	if token == nil || token.AccessToken == "" {
		return nil, errors.New("required token")
	}
	NewClient, err := NewClient(token, baseUrl, client)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &YandexDisk{Token: token, Client: NewClient}, nil
}
