package local

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"io"
	"mime/multipart"
	"os"
)

func (str *storage) Save(_ context.Context, apkFile multipart.File, nameFile string) (string, error) {
	str.Lock()

	defer str.Unlock()

	filePath := fmt.Sprintf("%s/%s", str.directory, nameFile)

	tempFile, err := os.Create(filePath)
	if err != nil {
		return "", errors.NewErrorf("Cannot create temp file: %s", "Не удалось создать временный файл: %s", err)
	}

	defer tempFile.Close()

	_, err = io.Copy(tempFile, apkFile)
	if err != nil {
		return "", errors.NewErrorf("Cannot copy to temp file: %s", "Не удалось скопировать во временный файл: %s", err)
	}

	if _, err := apkFile.Seek(0, io.SeekStart); err != nil {
		return "", errors.NewErrorf("Cannot set seek: %v", "Не удалось выставить начало файла: %v", err)
	}

	return filePath, nil
}
