package models

import "github.com/Ponchitos/application_service/server/tools/types"

type ApplicationTrack struct {
	ID                  int `json:"-"`
	GoogleApplicationID int `json:"-"`

	Created types.NullTime `json:"created"`

	TrackID    string `json:"trackId"`
	TrackAlias string `json:"trackAlias"`
}
