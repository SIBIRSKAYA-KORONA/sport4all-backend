package models

// swagger:model Sport
type Sport struct {
	// example: 10
	ID uint `json:"id" gorm:"primary_key"`

	// example: Football
	Name string `json:"name" gorm:"unique"`

	About string `json:"about"`

	// example: https://someurl
	Avatar Attach `json:"avatar" gorm:"foreignKey:sportName"`

	Tournament []Tournament `json:"tournament,omitempty" gorm:"foreignKey:sportName"`
}

func (sport *Sport) TableName() string {
	return "sport"
}
