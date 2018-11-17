package models

// Resource : resource field
type Resource struct {
	ID               int    `gorm:"primary_key"`
	APIID            int    `json:"api_id"`
	Path             string `json:"path" gorm:"size:2048;unique"`
	RegisteredMethod string `json:"registed_http_method" gorm:"size:64"`
	DestMethod       string `json:"dest_http_method" gorm:"size:64"`
	DestURL          string `json:"dest_url" gorm:"size:2048"`
	Timeout          int    `json:"timeout"`
}
