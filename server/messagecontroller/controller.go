package messagecontroller

import (
	applicationDecoder "github.com/Ponchitos/application_service/server/infrastructure/decoders/application"
	applicationEncoder "github.com/Ponchitos/application_service/server/infrastructure/encoders/application"
	"github.com/Ponchitos/application_service/server/internal/endpoints/application"
	"github.com/Ponchitos/application_service/server/internal/services/applications"
	"github.com/Ponchitos/application_service/server/pkg"
	"github.com/Ponchitos/application_service/server/pkg/mqcoordinator"
	coordinator2 "github.com/Ponchitos/application_service/server/pkg/mqcoordinator/coordinator"
	"github.com/Ponchitos/application_service/server/tools/logger"
)

func NewMessageController(
	logger logger.Logger,
	subscriber pkg.Subscriber,
	publisher pkg.Publisher,
	service applications.Service,
) (mqcoordinator.MsCoordinator, error) {

	coordinator, err := coordinator2.NewCoordinator(coordinator2.CoordinatorConfig{}, logger)
	if err != nil {
		return nil, err
	}

	changeApplicationStatusEndpoint := application.MakeChangeApplicationStatusEndpoint(service)
	uninstallApplicationCompleteEndpoint := application.MakeUninstallApplicationCompleteEndpoint(service)

	coordinator.AddHandler(
		"changeApplicationStatus",
		"application.change_application_status",
		subscriber,
		"application.change_application_status.complete",
		publisher,
		changeApplicationStatusEndpoint,
		applicationDecoder.ChangeApplicationStatusDecoder,
		applicationEncoder.ChangeApplicationStatusEncoder,
	)

	coordinator.AddHandler(
		"uninstallApplicationComplete",
		"application.uninstall_application_complete",
		subscriber,
		"application.uninstall_application_complete.complete",
		publisher,
		uninstallApplicationCompleteEndpoint,
		applicationDecoder.UninstallApplicationCompleteDecoder,
		applicationEncoder.UninstallApplicationCompleteEncoder,
	)

	return coordinator, nil
}
