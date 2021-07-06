package apk

import (
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/infrastructure/filestorages"
	"github.com/Ponchitos/application_service/server/infrastructure/filestorages/local"
	"github.com/Ponchitos/application_service/server/infrastructure/filestorages/s3"
	"github.com/Ponchitos/application_service/server/internal/handlers"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"os"
	"os/exec"
	"path"
)

const (
	scriptName = "./server/infrastructure/handlers/apk/target.py"
)

type apkHandler struct {
	config     *config.Config
	lgr        logger.Logger
	storages   map[string]filestorages.FileStorage
	pythonPath string
	scriptPath string
}

func NewApkHandler(config *config.Config, lgr logger.Logger) (handlers.ApkFileHandler, error) {
	handler := &apkHandler{config: config, lgr: lgr, storages: make(map[string]filestorages.FileStorage)}

	err := handler.initStorages()
	if err != nil {
		return nil, err
	}

	err = handler.initPythonScript()
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func (hAPK *apkHandler) initStorages() error {
	hAPK.storages[handlers.LocalMode] = local.NewLocalStorage(hAPK.config, hAPK.lgr)

	s3Storage, err := s3.NewS3Storage(hAPK.config, hAPK.lgr)
	if err != nil {
		return err
	}

	hAPK.storages[handlers.ExternalMode] = s3Storage

	return nil
}

func (hAPK *apkHandler) initPythonScript() error {
	pythonPath, err := exec.LookPath(hAPK.config.PythonName)
	if err != nil {
		return errors.NewErrorf("Cannot find %s: %v", "Не удалось найти команду %s: %v", hAPK.config.PythonName, err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		return errors.NewErrorf("Cannot find work dir: %v", "Не удалось определить рабочую папку: %v", err)
	}

	scriptPath := path.Join(workDir, scriptName)

	_, err = os.Stat(scriptPath)
	if err != nil {
		return errors.NewErrorf("Cannot check script file exist: %v", "Не удалось проверить существование скрипта: %v", err)
	}

	hAPK.pythonPath = pythonPath
	hAPK.scriptPath = scriptPath

	return nil
}

func (hAPK *apkHandler) getStore(mode string) (filestorages.FileStorage, error) {
	store, ok := hAPK.storages[mode]
	if !ok {
		return nil, errors.NewError("Mode not defined", "Режим не определен")
	}

	return store, nil
}
