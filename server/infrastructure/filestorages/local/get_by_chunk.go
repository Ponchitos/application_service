package local

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"os"
)

func (str *storage) GetByChunk(_ context.Context, fileName string, offset, whence int64) ([]byte, int64, error) {
	str.Lock()

	defer str.Unlock()

	filePath := fmt.Sprintf("%s/%s", str.directory, fileName)
	chunk := make([]byte, whence-offset)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, errors.NewErrorf("Cannot open file: %s", "Не удалось открыть файл: %s", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, errors.NewErrorf("Cannot get file info: %s", "Не удалось получить информацию о файле: %s", err)
	}

	_, err = file.Seek(offset, int(whence))
	if err != nil {
		return nil, 0, errors.NewErrorf("Cannot set seek: %s", "Не удалось выставить ограничения: %s", err)
	}

	_, err = file.Read(chunk)
	if err != nil {
		return nil, 0, errors.NewErrorf("Cannot read from file: %s", "Не удалось прочитать из файла: %s", err)
	}

	return chunk, fileInfo.Size(), nil
}
