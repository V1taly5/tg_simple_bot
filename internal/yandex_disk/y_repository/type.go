package yrepository

import (
	"net/url"
	"time"
)

// add yandexDisk interface

type YandexDisk struct {
	Token  *Token
	Client *Client
}

type Share struct {
	IsRoot  bool   `json:"is_root"`
	IsOwned bool   `json:"is_owned"`
	Rights  string `json:"rights"`
}

type Embedded struct {
	Sort   string     `json:"sort"`
	Items  []struct{} `json:"items"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
	Path   string     `json:"path"`
	Total  int        `json:"total"`
}

type Exif struct {
	DateTime     time.Time `json:"date_time"`
	GpsLongitude struct{}  `json:"gps_longitude"`
	GpsLatitude  struct{}  `json:"gps_latitude"`
}

type Size struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type CommentIDs struct {
	PrivateResource string `json:"private_resource"`
	PublicResource  string `json:"public_resource"`
}

// Response /v1/disk/resources Get MetaInfo about file or dir
type Resources struct {
	AntivirusStatus  map[string]interface{} `json:"antivirus_status"`
	ResourceID       string                 `json:"resource_id"`
	Share            Share                  `json:"share"`
	File             string                 `json:"file"`
	Size             int                    `json:"size"`
	PhotoSliceTime   time.Time              `json:"photoslice_time"`
	Embedded         Embedded               `json:"_embedded"`
	Exif             Exif                   `json:"exif"`
	CustomProperties map[string]interface{} `json:"custom_properties"`
	MediaType        string                 `json:"media_type"`
	Preview          string                 `json:"preview"`
	Type             string                 `json:"type"`
	MimeType         string                 `json:"mime_type"`
	Revision         int                    `json:"revision"`
	PublicURL        string                 `json:"public_url"`
	Path             string                 `json:"path"`
	MD5              string                 `json:"md5"`
	PublicKey        string                 `json:"public_key"`
	SHA256           string                 `json:"sha256"`
	Name             string                 `json:"name"`
	Created          time.Time              `json:"created"`
	Sizes            []Size                 `json:"sizes"`
	Modified         time.Time              `json:"modified"`
	CommentIDs       CommentIDs             `json:"comment_ids"`
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

func (e *Error) Is(target error) bool {
	// Проверяем, является ли target экземпляром Error
	targetErr, ok := target.(*Error)
	if !ok {
		return false
	}

	// Сравниваем поля ошибок
	if e.Message != targetErr.Message || e.Description != targetErr.Description || e.ErrorMSG != targetErr.ErrorMSG {
		return false
	}

	return true
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
