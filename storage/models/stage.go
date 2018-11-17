package models

// Stage : Stage field
type Stage struct {
	ID             int    `gorm:"primary_key"`
	APIID          int    `json:"api_id"`
	StageName      string `json:"stage_name"`
	StageVariables string `json:"stage_variables" sql:"type:text"`
}
