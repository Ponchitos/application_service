package local

import (
	"github.com/Ponchitos/application_service/server/config"
	"github.com/Ponchitos/application_service/server/infrastructure/filestorages"
	"github.com/Ponchitos/application_service/server/tools/logger"
	"os"
	"sync"
)

type storage struct {
	sync.Mutex

	config    *config.Config
	lgr       logger.Logger
	directory string
}

func NewLocalStorage(config *config.Config, lgr logger.Logger) filestorages.FileStorage {
	return &storage{config: config, lgr: lgr, directory: os.TempDir()}
}
