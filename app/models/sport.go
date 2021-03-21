package models

// swagger:model Sport
type Sport struct {
	// example: 10
	ID uint `json:"id" gorm:"primary_key"`

	// example: Football
	Kind string `json:"kind" gorm:"unique;index"`

	// example: https://someurl
	ThumbURL string `json:"thumb" gorm:"unique"`
}
