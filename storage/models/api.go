package models

// API : Base info for a group of resource
type API struct {
	ID          int    `gorm:"primary_key"`
	APIName     string `json:"api_name" gorm:"size:64;unique"`
	APINameSlug string `json:"api_name_slug" gorm:"size:64;unique"`
}
