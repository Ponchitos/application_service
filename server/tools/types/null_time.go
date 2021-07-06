package types

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const layout = "2006-01-02T15:04:05"

type NullTime struct {
	sql.NullTime
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf(`"%s"`, nt.Time.Format(layout))), nil
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		return nil
	}

	var err error
	nt.Time, err = time.Parse(layout, s)
	nt.Valid = err == nil

	return err
}
