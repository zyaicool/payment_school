package models

import "time"

type HistoryMasterData struct {
	Master
	GenerateDate time.Time `json:"generateDate"`
	FileName     string    `json:"fileName"`
}
