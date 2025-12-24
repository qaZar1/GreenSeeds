package models

type DeviceSettings struct {
	Key   string `json:"key" db:"key"`
	Value string `json:"value" db:"value"`
} // @name deviceSettings
