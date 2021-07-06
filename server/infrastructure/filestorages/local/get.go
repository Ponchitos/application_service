package local

import (
	"context"
	"fmt"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"io/ioutil"
)

func (str *storage) Get(_ context.Context, nameFile string) ([]byte, error) {
	str.Lock()

	defer str.Unlock()

	filePath := fmt.Sprintf("%s/%s", str.directory, nameFile)

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.NewErrorf("Cannot open file: %s", "Не удалось открыть файл: %s", err)
	}

	return file, nil
}
