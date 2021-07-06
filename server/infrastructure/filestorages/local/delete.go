package local

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"os"
)

func (str *storage) Delete(_ context.Context, filePath string) error {
	str.Lock()

	_, err := os.Stat(filePath)
	if err != nil {
		return errors.NewErrorf("Cannot check file exist: %s", "Не удалось проверить существование файла: %s", err)
	}

	err = os.Remove(filePath)
	if err != nil {
		return errors.NewErrorf("Cannot remove file: %v", "Не удалось удалить файл: %v", err)
	}

	str.Unlock()

	return nil
}
