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

func (d *YandexDisk) GetDisk(ctx context.Context, fields []string) (*Disk, error) {
	const op = "ydisk/GetDisk"
	//add checking for zero
	values := buildURLValues("fields", strings.Join(fields, ","))

	req, err := d.Client.MakeRequest(http.MethodGet, "/disk?"+values.Encode(), nil)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	disk := new(Disk)

	_, err = d.Client.GetResponse(ctx, req, &disk)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return disk, nil
}

func (d *YandexDisk) GetResourceUploadLink(
	ctx context.Context,
	path string,
	fields []string,
	overwrite bool) (
	*ResourceUploadLink,
	error) {
	const op = "ydisk/GetResourceUploadLink"

	values := buildURLValues("paht", path,
		"fields", strings.Join(fields, ","),
		"overwrite", strconv.FormatBool(overwrite),
	)

	req, err := d.Client.MakeRequest(http.MethodGet, "/disk/resources/upload?"+values.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resourceUploadLink := new(ResourceUploadLink)

	_, err = d.Client.GetResponse(ctx, req, &resourceUploadLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resourceUploadLink, nil
}

func (d *YandexDisk) UploudExternalResource(ctx context.Context, path string, externalUrl string, disableRedirects bool, fields []string) (*Link, error) {
	const op = "ydisk/UploudExternalResource"

	values := buildURLValues("path", path,
		"url", externalUrl,
		"disable_redirects", strconv.FormatBool(disableRedirects),
		"fields", strings.Join(fields, ","),
	)

	req, err := d.Client.MakeRequest(http.MethodPost, "/disk/resources/upload?"+values.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	link := new(Link)
	_, err = d.Client.GetResponse(ctx, req, &link)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return link, nil
}
func buildURLValues(pairs ...string) *url.Values {
	values := &url.Values{}
	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i]
		value := pairs[i+1]
		if value != "" {
			values.Add(key, value)
		}
	}
	return values
}
