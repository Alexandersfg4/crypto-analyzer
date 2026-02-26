package models

import "time"

type FearAndGreed struct {
	LastWeek struct {
		Timestamp           int    `json:"timestamp,omitzero"`
		Value               int    `json:"value,omitzero"`
		ValueClassification string `json:"value_classification,omitzero"`
	} `json:"lastWeek"`
	Name string `json:"name,omitzero"`
	Now  struct {
		Timestamp           int       `json:"timestamp,omitzero"`
		UpdateTime          time.Time `json:"update_time,omitzero"`
		Value               int       `json:"value,omitzero"`
		ValueClassification string    `json:"value_classification,omitzero"`
	} `json:"now"`
	Yesterday struct {
		Timestamp           int    `json:"timestamp,omitzero"`
		Value               int    `json:"value,omitzero"`
		ValueClassification string `json:"value_classification,omitzero"`
	} `json:"yesterday"`
}
