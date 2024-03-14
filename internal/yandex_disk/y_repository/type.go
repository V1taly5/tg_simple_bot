package yrepository

import "net/url"

// add yandexDisk interface

type YandexDisk struct {
	Token  *Token
	Client *Client
}

type SystemFolders struct {
	Odnoklassniki string `json:"odnoklassniki"`
	Google        string `json:"google"`
	Insagram      string `json:"insagram"`
	Vkontakte     string `json:"vkontakte"`
	Attach        string `json:"attach"`
	Mailru        string `json:"mailru"`
	Downloads     string `json:"downloads"`
	Applications  string `json:"applications"`
	Facebook      string `json:"facebook"`
	Social        string `json:"social"`
	Messenger     string `json:"messenger"`
	Calendar      string `json:"calendar"`
	Photostream   string `json:"photostream"`
	Screenshots   string `json:"screenshots"`
	Scans         string `json:"scans"`
}

type User struct {
	RegTime     string `json:"reg_time"`
	DisplayName string `json:"display_name"`
	Uid         string `json:"uid"`
	Country     string `json:"country"`
	IsChild     bool   `json:"is_child"`
	Login       string `json:"login"`
}

type Disk struct {
	PaidMaxFileSize            int    `json:"paid_max_file_size"`
	MaxFileSize                int    `json:"max_file_size"`
	TotalSpace                 int    `json:"total_space"`
	TrashSize                  int    `json:"trash_size"`
	UsedSpace                  int    `json:"used_space"`
	IsPaid                     bool   `json:"is_paid"`
	RegTime                    string `json:"reg_time"`
	SystemFolders              SystemFolders
	User                       User
	UnlimitedAutouploadEnabled bool `json:"unlimited_autoupload_enabled"`
	Revision                   int  `json:"revision"`
}

type Link struct {
	Href      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}

type ResourceUploadLink struct {
	OperationID string `json:"operation_id"`
	Href        string `json:"href"`
	Method      string `json:"method"`
	Templated   bool   `json:"templated"`
}

type Token struct {
	AccessToken string
}

type ResponseInfo struct {
	Status     string
	StatusCode int
}

func (r *ResponseInfo) SetResponseInfo(status string, statusCode int) {
	r.Status = status
	r.StatusCode = statusCode
}

type Error struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	ErrorMSG    string `json:"error"`
}

func (e *Error) Error() string {
	return e.ErrorMSG
}

type Params map[string]string

func buildParams(in Params) url.Values {
	if in == nil {
		return url.Values{}
	}
	out := url.Values{}
	for key, value := range in {
		out.Set(key, value)
	}
	return out
}
