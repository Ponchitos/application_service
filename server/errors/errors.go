package errors

import "github.com/Ponchitos/application_service/server/infrastructure/errors"

var (
	BadRequest   = errors.NewError("Not valid request", "Некорректный запрос")
	NotValidType = errors.NewError("Not valid type value", "Не валидный тип значения")
)
