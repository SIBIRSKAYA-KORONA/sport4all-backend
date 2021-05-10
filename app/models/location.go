package models

type Location struct {
	City        string `json:"city"`
	Country     string `json:"country" gorm:"not null, index"`
	Continent   string `json:"continent"`
	Subdivision string `json:"subdivisions" gorm:"index"`
}
