package application

import (
	"context"
	__ "github.com/Ponchitos/application_service/server/infrastructure/decoders/application/proto"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/endpoints/application"
	"google.golang.org/protobuf/proto"
)

func UninstallApplicationCompleteDecoder(ctx context.Context, message interface{}) (interface{}, error) {
	var receiveMessage __.UninstallApplicationCompleteRequest

	messageOfBytes, ok := message.([]byte)
	if !ok {
		return nil, errors.NewError("Received message not valid", "Полученное сообщение некорректно")
	}

	err := proto.Unmarshal(messageOfBytes, &receiveMessage)
	if err != nil {
		return nil, errors.NewErrorf("Can't unmarshal message: %v", "Не удалось десериализировать сообщение: %v", err)
	}

	return &application.UninstallApplicationCompleteRequest{
		VersionUUID:     receiveMessage.VersionUUID,
		EnterpriseID:    receiveMessage.EnterpriseID,
		ApplicationUUID: receiveMessage.ApplicationUUID,
	}, nil
}
