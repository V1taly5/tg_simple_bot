package yrepository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

type client struct {
	httpClient *http.Client
	token      *Token
	baseUrl    *url.URL
}

func NewClient(token *Token, baseUrl string, httpClient *http.Client) (*client, error) {
	if token == nil || token.AccessToken == "" {
		return nil, errors.New("required token")
	}
	if baseUrl == "" {
		return nil, errors.New("required baseUrl")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	base, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	c := &client{httpClient: httpClient, token: token, baseUrl: base}
	return c, nil
}

func (c *client) SetRequestHeaders(req *http.Request) {
	req.Header.Add("Access", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "OAuth "+c.token.AccessToken)
}

func (c *client) MakeRequest(method string, pathUrl string, body io.Reader) (*http.Request, error) {
	fullEndPoint := c.baseUrl.ResolveReference(&url.URL{Path: pathUrl})

	req, err := http.NewRequest(method, fullEndPoint.String(), body)
	if err != nil {
		return nil, err
	}

	c.SetRequestHeaders(req)

	return req, nil
}

func (c *client) DoRequset(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	return resp, err
}

func (c *client) GetResponse(ctx context.Context, req *http.Request, obj interface{}) (*ResponseInfo, error) {
	const op = "yadisk/GetResponde"
	resp, err := c.DoRequset(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() {
		log := slog.New(slog.NewTextHandler(os.Stdout, nil))
		if err := bodyClose(resp.Body); err != nil {
			log.Error("io.Close not closed")
		}
	}()
	responseInfo := new(ResponseInfo)
	responseInfo.SetResponseInfo(resp.Status, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseInfo, err
	}
	if len(body) > 0 {
		// Попробуем разобрать JSON в структуру Error
		errorJSON := new(Error)
		err := json.Unmarshal(body, &errorJSON)
		if err != nil {
			return responseInfo, fmt.Errorf("%s: %w", op, err)
		} else if (Error{} != *errorJSON) {
			// Проверяем пришла ли нам ошибка со стороннего API
			return responseInfo, fmt.Errorf("%s: %w", op, errorJSON)
		}
		// Если JSON не содержит информации об ошибке, попробуем разобрать в другую структуру (obj)
		err = json.Unmarshal(body, &obj)
		if err != nil {
			return responseInfo, fmt.Errorf("%s: %w", op, err)
		}
	}
	return responseInfo, nil
}

func bodyClose(closer io.Closer) error {
	err := closer.Close()
	if err != nil {
		return err
	}
	return nil
}
