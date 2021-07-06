package application

import (
	"context"
	"github.com/Ponchitos/application_service/server/infrastructure/decoders/application/proto"
	"github.com/Ponchitos/application_service/server/infrastructure/errors"
	"github.com/Ponchitos/application_service/server/internal/endpoints/application"
	"google.golang.org/protobuf/proto"
)

func ChangeApplicationStatusDecoder(ctx context.Context, message interface{}) (interface{}, error) {
	var receiveMessage __.ChangeApplicationStatusRequest

	messageOfBytes, ok := message.([]byte)
	if !ok {
		return nil, errors.NewError("Received message not valid", "Полученное сообщение некорректно")
	}

	err := proto.Unmarshal(messageOfBytes, &receiveMessage)
	if err != nil {
		return nil, errors.NewErrorf("Can't unmarshal message: %v", "Не удалось десериализировать сообщение: %v", err)
	}

	return &application.ChangeApplicationStatusRequest{
		VersionUUID:  receiveMessage.VersionUUID,
		EnterpriseID: receiveMessage.EnterpriseID,
		Status:       receiveMessage.Status.String(),
	}, nil
}
