package models

import ()

type RecordEnvelope struct {
	Type string `json:"_type"`
	Data Record `json:"data"`
}

type RecordListEnvelope struct {
	Records int      `json:"records"`
	Type    string   `json:"_type"`
	Data    []Record `json:"data"`
}

type WaterEntityEnvelope struct {
	Type string   `json:"_type"`
	Data []Record `json:"data"`
}

type SettingsEnvelope struct {
	Data UserSettings `json:"data"`
}

type LogEnvelope struct {
	Type string  `json:"_type"`
	Data UserLog `json:"data"`
}
