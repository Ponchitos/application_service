package models

import (
	"encoding/json"
	"github.com/Ponchitos/application_service/server/tools/types"
)

type ManagedProperty struct {
	GoogleApplicationID int `json:"-"`
	ManagedPropertyID   int `json:"id"`

	Created types.NullTime `json:"created"`

	Key         string `json:"key"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`

	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Value        interface{} `json:"value,omitempty"`

	Entries          []*ManagedPropertyEntry `json:"entries,omitempty"`
	NestedProperties []*ManagedProperty      `json:"nestedProperties,omitempty"`
}

func (property *ManagedProperty) ConvertEntriesToBytes() ([]byte, error) {
	return json.Marshal(property.Entries)
}

func (property *ManagedProperty) ConvertValueToBytes() ([]byte, error) {
	return json.Marshal(property.Value)
}

func (property *ManagedProperty) ConvertBytesToEntries(data []byte) error {
	if data == nil {
		property.Entries = make([]*ManagedPropertyEntry, 0)

		return nil
	}

	return json.Unmarshal(data, &property.Entries)
}

func (property *ManagedProperty) ConvertBytesToValue(data []byte) error {
	if data == nil {
		return nil
	}

	return json.Unmarshal(data, &property.Value)
}
