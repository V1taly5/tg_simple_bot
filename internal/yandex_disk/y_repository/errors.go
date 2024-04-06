package yrepository

var ErrResourceNotFound = &Error{
	Message:     "Не удалось найти запрошенный ресурс.",
	Description: "Resource not found.",
	ErrorMSG:    "DiskNotFoundError",
}
