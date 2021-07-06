package errors

import (
	"fmt"
	"os"
)

var useAltErrLang bool

func init() {
	value, ok := os.LookupEnv("ALT_ERR_LANG")
	if ok {
		switch value {
		case "RUS", "Rus", "rus", "RU", "Ru", "ru", "1", "true", "TRUE":
			useAltErrLang = true
		}
	}
}

type errorString struct {
	s    string
	altS string
}

func (e *errorString) Error() string {
	if useAltErrLang {
		return e.altS
	}

	return e.s
}

func NewError(text, altText string) error {
	return &errorString{
		s:    text,
		altS: altText,
	}
}

func NewErrorf(template, alttemplate string, value ...interface{}) error {
	return &errorString{
		s:    fmt.Sprintf(template, value...),
		altS: fmt.Sprintf(alttemplate, value...),
	}
}
