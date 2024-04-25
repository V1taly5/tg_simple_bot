package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	yrepository "simple_tg_bot/internal/yandex_disk/y_repository"
)

const (
	directory string = "dir"
)

var (
	ErrPathIsEmpty      = errors.New("path is empty")
	ErrPathNotDirectory = errors.New("end of the path is not a directory")
)

type DiskUseCase struct {
	diskRepository *yrepository.YandexDisk
}

func NewDiskUseCase(d *yrepository.YandexDisk) *DiskUseCase {
	return &DiskUseCase{d}
}

func (yr *DiskUseCase) UploudFileByURL(ctx context.Context,
	log *slog.Logger,
	pathToSave string,
	fileNameToSave string,
	fileDownloadURL string,
	disableRedirect bool,
	fields []string) (*yrepository.Link, error) {
	const op = "yandex_disk/usecase/UploadFileByURL"

	mas := []string{"type"}

	if pathToSave == "" {
		log.Debug("path is empty")
		return nil, fmt.Errorf("%s: %w", op, ErrPathIsEmpty)
	}
	res, err := yr.diskRepository.GetMetaInfoFile(ctx, pathToSave, mas, 10, 0, true, "100", "type")
	if err != nil {
		if errors.Is(err, yrepository.ErrResourceNotFound) {
			log.Warn("path not found")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if res.Type != directory {
		return nil, fmt.Errorf("%s: %w", op, ErrPathNotDirectory)
	}

	result, err := yr.diskRepository.UploudExternalResource(ctx, fmt.Sprintf("%s/%s", pathToSave, fileNameToSave), fileDownloadURL, disableRedirect, fields)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}
